package jwtx

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hanyougame/glib/metadata"
	"github.com/hanyougame/glib/utils"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
	"net/http"
	"time"
)

var (
	ErrInvalidToken     = errors.New("invalid auth token")
	ErrNoClaims         = errors.New("no auth params")
	ErrNotDetailField   = errors.New("no detail field")
	ErrJwtGenerateError = errors.New("token generate error")
)

const cacheJwtInfoKey = "{cache}:jwt:%s:%v"

type JWT struct {
	rdb    redis.UniversalClient //redis
	sso    bool                  // 单点登录验证redis 必须为true
	secret string                // JWT密钥
	ttl    int64                 // 过期时间
	scene  string                //场景值
}

func NewJWT() *JWT {
	return &JWT{}
}

func (j *JWT) WithRdb(rdb redis.UniversalClient) *JWT {
	j.rdb = rdb
	return j
}

func (j *JWT) WithSso(sso bool) *JWT {
	j.sso = sso
	return j
}

func (j *JWT) WithSecret(secret string) *JWT {
	j.secret = secret
	return j
}

func (j *JWT) WithTTL(ttl int64) *JWT {
	j.ttl = ttl
	return j
}

func (j *JWT) WithScene(scene string) *JWT {
	j.scene = scene
	return j
}

func (j *JWT) GenerateToken(ctx context.Context, uid any, claims jwt.MapClaims) (string, error) {
	now := time.Now().Unix()
	payload := make(jwt.MapClaims)
	payload[metadata.CtxJWTUserId] = uid
	if claims != nil {
		for k, v := range claims {
			payload[k] = v
		}
	}
	accessToken, err := j.getJwtToken(j.secret, now, j.ttl, payload)
	if err != nil {
		logx.Errorf("get jwt token err:%v", err)
		return "", ErrJwtGenerateError
	}
	if j.sso && j.rdb != nil {
		err = j.rdb.Set(ctx, j.generateCacheKey(uid), accessToken, time.Second*time.Duration(j.ttl)).Err()
		if err != nil {
			logx.Error("set jwt info to redis err: ", err)
			return "", ErrJwtGenerateError
		}
	}
	return accessToken, nil
}

func (j *JWT) getJwtToken(secretKey string, iat, seconds int64, payload jwt.MapClaims) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payload {
		claims[k] = v
	}
	token2 := jwt.New(jwt.SigningMethodHS256)
	token2.Claims = claims
	return token2.SignedString([]byte(secretKey))
}

func (j *JWT) ParseToken(r *http.Request) (uid any, claims jwt.MapClaims, err error) {
	var (
		tok    *jwt.Token
		parser = token.NewTokenParser()
		ok     bool
	)

	if tok, err = parser.ParseToken(r, j.secret, j.secret); err != nil {
		err = ErrInvalidToken
		return
	}

	if !tok.Valid {
		err = ErrInvalidToken
		return
	}

	if claims, ok = tok.Claims.(jwt.MapClaims); !ok {
		err = ErrNoClaims
		return
	}

	if uid, ok = claims[metadata.CtxJWTUserId]; !ok {
		err = ErrNotDetailField
		return
	}

	if j.sso && j.rdb != nil { //验证token是否存在与redis
		var stringCmd = j.rdb.Get(r.Context(), j.generateCacheKey(uid))
		if err = stringCmd.Err(); err != nil {
			if errors.Is(err, redis.Nil) {
				err = ErrInvalidToken
				return
			}
			return
		}
		if stringCmd.Val() != tok.Raw {
			err = ErrInvalidToken
			return
		}
	}

	return
}

func (j *JWT) generateCacheKey(uid any) string {
	return fmt.Sprintf(cacheJwtInfoKey, utils.Ternary(j.scene == "", j.secret, j.scene), uid)
}

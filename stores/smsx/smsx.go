package sms

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/hanyougame/glib/utils/httpc"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/jsonx"
	"time"
)

type (
	Sms struct {
		Config Config `json:"config"`
	}

	Config struct {
		URL       string `json:"url"`
		AppCode   string `json:"app_code"`
		AppKey    string `json:"app_key"`
		AppSecret string `json:"app_secret"`
	}

	request struct {
		AppCode   string `json:"appcode"`
		AppKey    string `json:"appkey"`
		Sign      string `json:"sign"`
		UID       string `json:"uid"`
		Phone     string `json:"phone"`
		Msg       string `json:"msg"`
		Timestamp int64  `json:"timestamp"`
	}
	Response struct {
		Code   string `json:"code"`
		Desc   string `json:"desc"`
		UID    string `json:"UID"`
		Result []struct {
			Status string `json:"status"`
			Phone  string `json:"phone"`
			Desc   string `json:"desc"`
		} `json:"result"`
	}
)

func NewSms(conf Config) *Sms {
	return &Sms{Config: conf}
}

func (s *Sms) Send(ctx context.Context, phone, msg string, uid ...int64) (result Response, err error) {
	r := request{
		AppCode:   s.Config.AppCode,
		AppKey:    s.Config.AppKey,
		Phone:     phone,
		Msg:       msg,
		Timestamp: time.Now().UnixMilli(),
	}
	if len(uid) > 0 {
		r.UID = cast.ToString(uid[0])
	}
	r.Sign = s.genSign(fmt.Sprintf("%s%s%d", r.AppKey, s.Config.AppSecret, r.Timestamp))
	resp, err := httpc.Do(ctx).SetBody(r).Post(s.Config.URL)
	if err != nil {
		return
	}

	if err = jsonx.Unmarshal(resp.Body(), &result); err != nil {
		return
	}

	return
}

func (s *Sms) genSign(data string) string {
	bt := md5.Sum([]byte(data))
	return hex.EncodeToString(bt[:])
}

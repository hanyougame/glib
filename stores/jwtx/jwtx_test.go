package jwtx

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"testing"
)

var j = NewJWT().
	WithScene("api").
	WithSecret("123456").
	WithSso(true).
	WithTTL(5)

func TestGenToken(t *testing.T) {

	token, err := j.GenerateToken(context.Background(), 1, jwt.MapClaims{
		"name": "1",
		"age":  1,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
}

func TestParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZ2UiOjEsImV4cCI6MTczOTAwMTk5OCwiaWF0IjoxNzM5MDAxOTkzLCJuYW1lIjoiMSIsInVpZCI6MX0.on-Ie0zcFwTYTnUq5BspTYbdT9BY9N4cZW0fHwA0M9E"

	data, err := j.ParseToken(&http.Request{
		Header: http.Header{"Authorization": []string{"Bearer " + token}},
	})
	t.Log(err)
	t.Log(data)
}

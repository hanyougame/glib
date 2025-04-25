package codec

import (
	"crypto/hmac"
	sha2562 "crypto/sha256"
	"encoding/hex"
)

var Hmac256 hmac256

type hmac256 struct {
}

func (hmac256) Encode(str, secret string) string {
	mac := hmac.New(sha2562.New, []byte(secret))
	_, _ = mac.Write([]byte(str))
	return hex.EncodeToString(mac.Sum(nil))
}

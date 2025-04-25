package codec

import (
	"crypto/sha512"
	"encoding/hex"
)

var Hmac512 hmac512

type hmac512 struct {
}

func (hmac512) Encode(str string) string {
	message := []byte(str)
	hash := sha512.New() //SHA-512加密
	hash.Write(message)
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}

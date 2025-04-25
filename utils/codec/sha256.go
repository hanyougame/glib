package codec

import (
	sha2562 "crypto/sha256"
	"encoding/hex"
)

var SHA256 sha256

type sha256 struct {
}

func (sha256) Encode(str string) string {
	hash := sha2562.New()
	hash.Write([]byte(str))

	return hex.EncodeToString(hash.Sum(nil))
}

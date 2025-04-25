package codec

import (
	"crypto/md5"
	"encoding/hex"
)

var MD5 md5Hash

type md5Hash struct {
}

func (md5Hash) Encode(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	hashedBytes := hash.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

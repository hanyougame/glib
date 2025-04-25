package codec

import "encoding/base64"

var Base64 base64Hash

type base64Hash struct {
}

func (base64Hash) UrlEncode(src string) string {
	return base64.URLEncoding.EncodeToString([]byte(src))
}

func (base64Hash) UrlDecode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

func (base64Hash) Encode(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

func (base64Hash) Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

package encoding

import "encoding/base64"

func EncodeToBase64(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

func DecodeBase64(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

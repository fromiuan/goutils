package crypto

import (
	"encoding/base64"
)

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table)

func EncodeBase64(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func DecodeBase64(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}

func DecodeBase64String(src string) (string, error) {
	bytes, err := coder.DecodeString(src)
	return string(bytes), err
}

func DecodeStdBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func EncodeStdBase64(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

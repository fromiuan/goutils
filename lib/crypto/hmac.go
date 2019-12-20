package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func HmacSignature(toSign, secret string) string {
	return hex.EncodeToString(HmacBytes([]byte(toSign), []byte(secret)))
}

func HmacBytes(toSign, secret []byte) []byte {
	_authSignature := hmac.New(sha256.New, secret)
	_authSignature.Write(toSign)
	return _authSignature.Sum(nil)
}

func HmacBase64Bytes(toSign, secret []byte) []byte {
	_authSignature := hmac.New(md5.New, secret)
	_authSignature.Write(toSign)
	return EncodeBase64(_authSignature.Sum(nil))
}

func CheckSignature(result, secret string, body []byte) bool {
	expected := HmacBytes(body, []byte(secret))
	resultBytes, err := hex.DecodeString(result)

	if err != nil {
		return false
	}
	return hmac.Equal(expected, resultBytes)
}

func Md5Signature(body []byte) string {
	_bodyMD5 := md5.New()
	_bodyMD5.Write([]byte(body))
	return hex.EncodeToString(_bodyMD5.Sum(nil))
}

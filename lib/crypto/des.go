package crypto

import (
	"crypto/cipher"
	"crypto/des"
)

func DesBase64Encrypt(dst, key, vector []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	size := block.BlockSize()
	dst = PKCS5Padding(dst, size)
	mode := cipher.NewCBCEncrypter(block, vector)
	crypted := make([]byte, len(dst))
	mode.CryptBlocks(crypted, dst)
	base64 := EncodeBase64(crypted)

	return base64, nil
}

func DesBase64Decrypt(dst, key, vector []byte) ([]byte, error) {
	base64, err := DecodeBase64(dst)
	if err != nil {
		return nil, err
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, vector)
	data := make([]byte, len(base64))
	mode.CryptBlocks(data, base64)
	data = PKCS5UnPadding(data)
	return data, nil
}

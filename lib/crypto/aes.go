package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// AesEncrypt 加密
func AesEncrypt(dst, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	size := block.BlockSize()
	dst = PKCS5Padding(dst, size)

	mode := cipher.NewCBCEncrypter(block, key[:size])
	crypted := make([]byte, len(dst))
	mode.CryptBlocks(crypted, dst)
	return crypted, nil
}

// AesDecrypt 解密
func AesDecrypt(dst, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	size := block.BlockSize()
	mode := cipher.NewCBCDecrypter(block, key[:size])
	data := make([]byte, len(dst))
	mode.CryptBlocks(data, dst)
	data = PKCS5UnPadding(data)
	return data, nil
}

func PKCS5Padding(ciphertext []byte, size int) []byte {
	padding := size - len(ciphertext)%size
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(dst []byte) []byte {
	length := len(dst)
	unpadding := int(dst[length-1])
	return dst[:(length - unpadding)]
}

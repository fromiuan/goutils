package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const (
	RSA_PRIVATE_KEY = iota
	RSA_PUBLIC_KEY
	RSA_CERT
)

type SHAwithRSA struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewSHAwithRSA() *SHAwithRSA {
	return &SHAwithRSA{}
}

// 签名
func (this *SHAwithRSA) Sign(data string, key []byte) ([]byte, error) {
	h := sha1.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	privateInterface, err := ParseKey(key, RSA_PRIVATE_KEY)
	if err != nil {
		return nil, err
	}
	privateKey := privateInterface.(*rsa.PrivateKey)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey,
		crypto.SHA1, hashed)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// 验证
func (this *SHAwithRSA) Verify(publicKey *rsa.PublicKey, data string, sig []byte) error {
	h := sha1.New()
	h.Write([]byte(data))
	hash := h.Sum(nil)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hash, sig)
}

// ----------------------------------------------------------------------------
func ParseKey(key []byte, kind int) (interface{}, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("pem decode key error")
	}
	var result interface{}
	var err error
	switch kind {
	case RSA_PRIVATE_KEY:
		result, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	case RSA_PUBLIC_KEY:
		result, err = x509.ParsePKIXPublicKey(block.Bytes)
	case RSA_CERT:
		var cert *x509.Certificate
		cert, err = x509.ParseCertificate(block.Bytes)
		if cert != nil && err == nil {
			result = cert.PublicKey
		}
	default:
		break
	}

	return result, err
}

// 加密
func Encrypt(origData []byte, key []byte) ([]byte, error) {

	pubInterface, err := ParseKey(key, RSA_PUBLIC_KEY)
	if err != nil {
		return nil, err
	}
	publicKey := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, origData)
}

// 解密
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {

	privateInterface, err := ParseKey(key, RSA_PRIVATE_KEY)
	if err != nil {
		return nil, err
	}
	privateKey := privateInterface.(*rsa.PrivateKey)
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

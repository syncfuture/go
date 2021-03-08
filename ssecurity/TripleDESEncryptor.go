package ssecurity

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

// CreateTripleDESEncryptor create triple des encryptor
func CreateTripleDESEncryptor(key string) *TripleDESEncryptor {
	return &TripleDESEncryptor{Key: key}
}

// TripleDESEncryptor _
type TripleDESEncryptor struct {
	Key string
}

// EncryptString encrypt a string
func (x *TripleDESEncryptor) EncryptString(in string) (string, error) {
	data := []byte(in)
	data, err := x.Encrypt(data)
	str := base64.StdEncoding.EncodeToString(data)
	return str, err
}

// DecryptString decrypt a string
func (x *TripleDESEncryptor) DecryptString(in string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return in, err
	}

	data, err = x.Decrypt(data)
	if err != nil {
		return in, err
	}

	return string(data), err
}

// Encrypt encrypt data
func (x *TripleDESEncryptor) Encrypt(in []byte) ([]byte, error) {
	key := []byte(x.Key)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	in = pkcs5Padding(in, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(in))
	blockMode.CryptBlocks(crypted, in)
	return crypted, nil
}

// Decrypt decrypt data
func (x *TripleDESEncryptor) Decrypt(in []byte) ([]byte, error) {
	key := []byte(x.Key)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(in))
	// origData := crypted
	blockMode.CryptBlocks(origData, in)
	origData = pkcs5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

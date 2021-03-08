package ssecurity

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"os"

	"github.com/syncfuture/go/srsautil"
)

// RSAEncryptor RSA encryptor
type RSAEncryptor struct {
	Key *rsa.PrivateKey
}

// CreateRSAEncryptorFromFile create encryptor by specifying cert path
func CreateRSAEncryptorFromFile(certPath string) (*RSAEncryptor, error) {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	return CreateRSAEncryptor(&data)
}

// CreateRSAEncryptor create RSA encryptor
func CreateRSAEncryptor(keyData *[]byte) (*RSAEncryptor, error) {
	r := RSAEncryptor{}
	var err error
	r.Key, err = srsautil.PKCS8BytesToPrivateKey(*keyData)
	return &r, err
}

// Encrypt encrypt
func (x *RSAEncryptor) Encrypt(in []byte) ([]byte, error) {
	data, err := rsa.EncryptPKCS1v15(rand.Reader, &x.Key.PublicKey, []byte(in))
	if err != nil {
		return in, err
	}

	return data, err
}

// Decrypt decrypt
func (x *RSAEncryptor) Decrypt(in []byte) ([]byte, error) {
	data, err := rsa.DecryptPKCS1v15(rand.Reader, x.Key, []byte(in))
	if err != nil {
		return in, err
	}
	return data, err
}

// EncryptString encrypt string
func (x *RSAEncryptor) EncryptString(in string) (string, error) {
	data, err := x.Encrypt([]byte(in))
	if err != nil {
		return in, err
	}

	// str := fmt.Sprintf("%x", data)
	str := base64.StdEncoding.EncodeToString(data)
	return str, err
}

// DecryptString decrypt string
func (x *RSAEncryptor) DecryptString(in string) (string, error) {
	inData, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return in, err
	}
	data, err := x.Decrypt(inData)
	if err != nil {
		return in, err
	}
	// str := fmt.Sprintf("%x", data)
	return string(data), err
}

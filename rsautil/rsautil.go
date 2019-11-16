package rsautil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"

	log "github.com/kataras/golog"
	u "github.com/syncfuture/go/util"
)

// GenerateKey generates a new key
func GenerateKey(bits int) (*rsa.PrivateKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if u.LogError(err) {
		return nil, err
	}
	return privkey, err
}

// PKCS1PrivateKeyToBytes private key to bytes
func PKCS1PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if u.LogError(err) {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, err
}

// PKCS1BytesToPrivateKey bytes to private key
func PKCS1BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if u.LogError(err) {
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	u.LogError(err)
	return key, err
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if u.LogError(err) {
			return nil, err
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if u.LogError(err) {
		return nil, err
	}

	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Error("not ok")
	}
	return key, err
}

// PKCS8PrivateKeyToBytes private key to bytes
func PKCS8PrivateKeyToBytes(priv *rsa.PrivateKey) ([]byte, error) {
	bytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if u.LogError(err) {
		return nil, err
	}

	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: bytes,
		},
	)

	return privBytes, err
}

// PKCS8BytesToPrivateKey bytes to private key
func PKCS8BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if u.LogError(err) {
			return nil, err
		}
	}
	key, err := x509.ParsePKCS8PrivateKey(b)
	u.LogError(err)
	return key.(*rsa.PrivateKey), err
}

// CertificateBytesToPublicKey bytes to public key
func CertificateBytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if u.LogError(err) {
			return nil, err
		}
	}
	ifc, err := x509.ParseCertificate(b)
	if u.LogError(err) {
		return nil, err
	}

	return ifc.PublicKey.(*rsa.PublicKey), err
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if u.LogError(err) {
		return nil, err
	}
	return ciphertext, err
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if u.LogError(err) {
		return nil, err
	}
	return plaintext, err
}

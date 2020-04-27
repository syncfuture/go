package security

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	log "github.com/syncfuture/go/slog"
	"golang.org/x/crypto/scrypt"
)

func GeneratePasswordSalt(keyLen int) string {
	salt := make([]byte, keyLen)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	saltString := base64.StdEncoding.EncodeToString([]byte(salt))
	return saltString
}

func HashPassword(salt, pass string, keyLen int) string {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		log.Error(err)
		return err.Error()
	}

	hash, err := scrypt.Key([]byte(pass), saltBytes, 1<<14, 8, 1, keyLen)
	if err != nil {
		log.Error(err)
		return err.Error()
	}

	hashString := base64.StdEncoding.EncodeToString([]byte(hash))
	return hashString
}

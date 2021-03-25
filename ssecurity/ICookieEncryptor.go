package ssecurity

import (
	"github.com/gorilla/securecookie"
	"github.com/syncfuture/go/serr"
)

type ICookieEncryptor interface {
	Encrypt(name string, value interface{}) (string, error)
	Decrypt(name, value string, dst interface{}) error
}

func NewSecureCookieEncryptor(s *securecookie.SecureCookie) ICookieEncryptor {
	return &SecureCookieEncryptor{
		s: s,
	}
}

type SecureCookieEncryptor struct {
	s *securecookie.SecureCookie
}

func (x *SecureCookieEncryptor) Encrypt(name string, value interface{}) (string, error) {
	a, b := x.s.Encode(name, value)
	return a, serr.WithStack(b)
}

func (x *SecureCookieEncryptor) Decrypt(name, value string, dst interface{}) error {
	err := x.s.Decode(name, value, dst)
	return serr.WithStack(err)
}

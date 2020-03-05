package security

import (
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12/context"
	"github.com/syncfuture/go/u"
)

type ISecureCookie interface {
	Set(ctx context.Context, cookieName, cookieValue string, options ...context.CookieOption) error
	Get(ctx context.Context, cookieName string) (string, error)
	Encode(name string, value interface{}) (string, error)
	Decode(name, value string, dst interface{}) error
}

type defaultSecureCookie struct {
	secure *securecookie.SecureCookie
}

func NewDefaultSecureCookie(hashKey, blockKey []byte) ISecureCookie {
	r := new(defaultSecureCookie)
	r.secure = securecookie.New(hashKey, blockKey)
	return r
}
func (x *defaultSecureCookie) Encode(name string, value interface{}) (string, error) {
	return x.secure.Encode(name, value)
}

func (x *defaultSecureCookie) Decode(name, value string, dst interface{}) error {
	return x.secure.Decode(name, value, dst)
}

func (x *defaultSecureCookie) Set(ctx context.Context, cookieName, cookieValue string, options ...context.CookieOption) error {
	encoded, err := x.secure.Encode(cookieName, cookieValue)

	if !u.LogError(err) {
		ctx.SetCookie(&http.Cookie{
			Name:     cookieName,
			Value:    encoded,
			HttpOnly: true,
			Secure:   true,
		}, options...)
	}

	return err
}

func (x *defaultSecureCookie) Get(ctx context.Context, cookieName string) (string, error) {
	cookieValue := ctx.GetCookie(cookieName)

	var decoded string
	err := x.secure.Decode(cookieName, cookieValue, &decoded)
	if u.LogError(err) {
		return "", err
	}

	return decoded, nil
}

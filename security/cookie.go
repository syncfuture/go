package security

import (
	"net/http"

	"github.com/gorilla/securecookie"
	log "github.com/kataras/golog"
	"github.com/kataras/iris/v12/context"
)

type ISecureCookie interface {
	Set(ctx context.Context, cookieName, cookieValue string, options ...context.CookieOption)
	Get(ctx context.Context, cookieName string) string
}

type defaultSecureCookie struct {
	secure *securecookie.SecureCookie
}

func NewDefaultSecureCookie(hashKey, blockKey []byte) ISecureCookie {
	r := new(defaultSecureCookie)
	r.secure = securecookie.New(hashKey, blockKey)
	return r
}

func (x *defaultSecureCookie) Set(ctx context.Context, cookieName, cookieValue string, options ...context.CookieOption) {
	if encoded, err := x.secure.Encode(cookieName, cookieValue); err == nil {
		ctx.SetCookie(&http.Cookie{
			Name:     cookieName,
			Value:    encoded,
			HttpOnly: true,
			Secure:   true,
		}, options...)
	} else {
		log.Error(err)
	}
}

func (x *defaultSecureCookie) Get(ctx context.Context, cookieName string) string {
	var decoded string
	cookieValue := ctx.GetCookie(cookieName)
	if err := x.secure.Decode(cookieName, cookieValue, &decoded); err == nil {
		log.Error(err)
	}

	return decoded
}

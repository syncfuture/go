package soidc

import (
	"github.com/kataras/iris/v12/context"
	"github.com/syncfuture/go/json"
	"github.com/syncfuture/go/security"
	u "github.com/syncfuture/go/util"
	"golang.org/x/oauth2"
)

const (
	key = "ecp:TOKENS:"
)

type CookieTokenStore struct {
	ctx          context.Context
	secureCookie security.ISecureCookie
}

func NewCookieTokenStore(ctx context.Context, secureCookie security.ISecureCookie) ITokenStore {
	r := new(CookieTokenStore)

	r.ctx = ctx
	r.secureCookie = secureCookie

	return r
}

func (x *CookieTokenStore) Save(userID string, token *oauth2.Token) error {
	j, err := json.Serialize(token)
	if u.LogError(err) {
		return err
	}

	err = x.secureCookie.Set(x.ctx, COKI_TOKEN, j)
	return err
}

func (x *CookieTokenStore) Get(userID string) (*oauth2.Token, error) {
	j, err := x.secureCookie.Get(x.ctx, COKI_TOKEN)
	if u.LogError(err) {
		return nil, err
	}

	t := new(oauth2.Token)
	err = json.Deserialize(j, t)
	u.LogError(err)
	return t, err
}

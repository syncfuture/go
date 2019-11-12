package soidc

import "golang.org/x/oauth2"

type ITokenStore interface {
	Save(userID string, token *oauth2.Token) error
	Get(userID string) (*oauth2.Token, error)
}

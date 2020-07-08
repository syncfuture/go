Tokenpackage soauth2

import "golang.org/x/oauth2"

type ITokenStore interface {
	GetToken(args ...interface{}) (*oauth2.Token, error)
	SaveToken(*oauth2.Token token) error
}

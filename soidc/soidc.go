package soidc

import (
	"net/http"

	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"github.com/syncfuture/go/security"
	"golang.org/x/oauth2"
)

const (
	SESS_ID       = "ID"
	SESS_USERNAME = "Username"
	SESS_EMAIL    = "Email"
	SESS_ROLES    = "Roles"
	SESS_LEVEL    = "Level"
	SESS_STATUS   = "Status"
	SESS_STATE    = "State"
	COKI_TOKEN    = ".ART"
	COKI_SESSION  = ".USS"
)

type IOIDCClient interface {
	HandleAuthentication(context.Context)
	HandleSignInCallback(context.Context)
	HandleSignOutCallback(context.Context)
	NewHttpClient(context.Context) (*http.Client, error)
	GetToken(ctx context.Context) (*oauth2.Token, error)
	SaveToken(ctx context.Context, token *oauth2.Token) error
}

type OIDCConfigs struct {
	ProjectName     string
	ClientID        string
	ClientSecret    string
	ProviderUrl     string
	CallbackURL     string
	AccessDeniedURL string
	Scopes          []string
}

type ClientOptions struct {
	ProviderUrl       string
	ClientID          string
	ClientSecret      string
	CallbackURL       string
	AccessDeniedURL   string
	Sess_ID           string
	Sess_Username     string
	Sess_Email        string
	Sess_Roles        string
	Sess_Level        string
	Sess_Status       string
	Sess_State        string
	Coki_Token        string
	Coki_Session      string
	Scopes            []string
	Sessions          *sessions.Sessions
	SecureCookie      security.ISecureCookie
	PermissionAuditor security.IPermissionAuditor
}

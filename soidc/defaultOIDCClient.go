package soidc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"

	"github.com/syncfuture/go/sslice"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/syncfuture/go/sjson"
	"github.com/syncfuture/go/u"

	"github.com/syncfuture/go/srand"

	oidc "github.com/coreos/go-oidc"
	log "github.com/kataras/golog"
	"github.com/kataras/iris/v12/context"

	gocontext "context"

	"golang.org/x/oauth2"
)

type defaultOIDCClient struct {
	Options      *ClientOptions
	OIDCProvider *oidc.Provider
	OAuth2Config *oauth2.Config
	ClientConfig *clientcredentials.Config
}

const offline_access = "offline_access"

func NewOIDCClient(options *ClientOptions) IOIDCClient {
	checkOptions(options)

	x := new(defaultOIDCClient)
	x.Options = options

	ctx := gocontext.Background()
	var err error
	x.OIDCProvider, err = oidc.NewProvider(ctx, options.PassportURL)
	if err != nil {
		log.Fatal(err)
	}
	x.OAuth2Config = new(oauth2.Config)
	x.OAuth2Config.ClientID = options.ClientID
	x.OAuth2Config.ClientSecret = options.ClientSecret
	x.OAuth2Config.Endpoint = x.OIDCProvider.Endpoint()
	x.OAuth2Config.RedirectURL = options.SignInCallbackURL
	x.OAuth2Config.Scopes = sslice.AppendStringToNew(options.Scopes, oidc.ScopeOpenID)

	x.ClientConfig = new(clientcredentials.Config)
	x.ClientConfig.ClientID = options.ClientID
	x.ClientConfig.ClientSecret = options.ClientSecret
	x.ClientConfig.TokenURL = x.OIDCProvider.Endpoint().TokenURL
	x.ClientConfig.Scopes = sslice.RemoveStringToNew(options.Scopes, oidc.ScopeOfflineAccess, oidc.ScopeOpenID)

	return x
}

func (x *defaultOIDCClient) NewHttpClient(args ...interface{}) (*http.Client, error) {
	goctx := gocontext.Background()
	if len(args) == 0 {
		// ClientCredential令牌方式
		return x.ClientConfig.Client(goctx), nil
	}

	ctx, ok := args[0].(context.Context)
	if !ok {
		panic("first parameter must be iris context")
	}

	session := x.Options.SessionManager.Start(ctx)
	userID := session.GetString(SESS_ID)
	if userID == "" {
		return http.DefaultClient, fmt.Errorf("user id not exists in session")
	}

	t, _, err := x.GetToken(ctx)
	if u.LogError(err) {
		return http.DefaultClient, err
	}

	tokenSource := x.OAuth2Config.TokenSource(goctx, t)
	newToken, err := tokenSource.Token()
	if u.LogError(err) {
		return http.DefaultClient, err
	}

	if newToken.AccessToken != t.AccessToken {
		x.SaveToken(ctx, newToken)
		log.Debugf("Saved new token for user %s", userID)
	}

	return oauth2.NewClient(goctx, tokenSource), nil
}

func (x *defaultOIDCClient) HandleAuthentication(ctx context.Context) {
	session := x.Options.SessionManager.Start(ctx)

	handlerName := ctx.GetCurrentRoute().MainHandlerName()
	area, controller, action := getRoutes(handlerName)

	// 判断请求是否允许访问
	userid := session.GetString(x.Options.Sess_ID)
	if userid != "" {
		roles := session.GetInt64Default(x.Options.Sess_Roles, 0)
		level := int32(session.GetIntDefault(x.Options.Sess_Level, 0))
		// 已登录
		allow := x.Options.PermissionAuditor.CheckRouteWithLevel(area, controller, action, roles, level)
		if allow {
			// 有权限
			ctx.Next()
			return

		} else {
			// 没权限
			ctx.Redirect(x.Options.AccessDeniedURL, http.StatusFound)
			return
		}
	} else {
		// 未登录
		allow := x.Options.PermissionAuditor.CheckRouteWithLevel(area, controller, action, 0, 0)
		if allow {
			// 允许匿名
			ctx.Next()
			return
		}

		// 记录请求地址，跳转去登录页面
		state := srand.String(32)
		session.Set(state, ctx.Request().URL.String())
		ctx.Redirect(x.OAuth2Config.AuthCodeURL(state), http.StatusFound)
	}
}

func (x *defaultOIDCClient) HandleSignIn(ctx context.Context) {
	returnURL := ctx.FormValue("returnUrl")
	if returnURL == "" {
		returnURL = "/"
	}

	session := x.Options.SessionManager.Start(ctx)
	userid := session.GetString(x.Options.Sess_ID)
	if userid != "" {
		// 已登录
		ctx.Redirect(returnURL, http.StatusFound)
		return
	}

	// 记录请求地址，跳转去登录页面
	state := srand.String(32)
	session.Set(state, returnURL)
	ctx.Redirect(x.OAuth2Config.AuthCodeURL(state), http.StatusFound)
}

func (x *defaultOIDCClient) HandleSignInCallback(ctx context.Context) {
	session := x.Options.SessionManager.Start(ctx)

	state := ctx.FormValue("state")
	redirectUrl := session.GetString(state)
	if redirectUrl == "" {
		ctx.WriteString("invalid state")
		ctx.StatusCode(http.StatusBadRequest)
		return
	}
	session.Delete(state) // 释放内存

	// 交换令牌
	code := ctx.FormValue("code")
	httpCtx := gocontext.Background()
	oauth2Token, err := x.OAuth2Config.Exchange(httpCtx, code /*,oauth2.SetAuthURLParam("aaa","bbb")*/)
	if err != nil {
		errStr := "Failed to exchange token: " + err.Error()
		log.Error(errStr)
		ctx.Write([]byte(errStr))
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	// 将字符串转化为令牌对象，忽略KeyFunc不存在错误
	jwtToken, err := new(jwt.Parser).Parse(oauth2Token.AccessToken, nil)
	vErr := err.(*jwt.ValidationError)
	if vErr.Errors != jwt.ValidationErrorUnverifiable {
		u.LogError(err)
		return
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok {
		session.Set(x.Options.Sess_ID, claims["sub"])
		session.Set(x.Options.Sess_Username, claims["name"])
		session.Set(x.Options.Sess_Roles, claims["role"])
		session.Set(x.Options.Sess_Level, claims["level"])
		session.Set(x.Options.Sess_Status, claims["status"])
		session.Set(x.Options.Sess_Email, claims["email"])

		// // 保存令牌
		x.SaveToken(ctx, oauth2Token)

		// 重定向到登录前页面
		ctx.Redirect(redirectUrl, http.StatusFound)
	} else {
		log.Error("cannot convert jwtToken.Claims to jwt.MapClaims")
	}
}

func (x *defaultOIDCClient) HandleSignOut(ctx context.Context) {
	session := x.Options.SessionManager.Start(ctx)
	_, idToken, _ := x.GetToken(ctx)

	session.Delete(SESS_ID)
	session.Delete(SESS_USERNAME)
	session.Delete(SESS_EMAIL)
	session.Delete(SESS_ROLES)
	session.Delete(SESS_LEVEL)
	session.Delete(SESS_STATUS)

	ctx.RemoveCookie(x.Options.Coki_Token)
	ctx.RemoveCookie(x.Options.Coki_IDToken)

	// 去Passport注销
	state := srand.String(32)
	session.Set(state, ctx.FormValue("returnUrl"))
	signoutUrl, _ := u.JointURLString(x.Options.PassportURL, "/connect/endsession")
	signoutUrl += "?post_logout_redirect_uri=" + url.PathEscape(x.Options.SignOutCallbackURL) + "&id_token_hint=" + idToken + "&state=" + state
	ctx.Redirect(signoutUrl, http.StatusFound)
}

func (x *defaultOIDCClient) HandleSignOutCallback(ctx context.Context) {
	session := x.Options.SessionManager.Start(ctx)

	state := ctx.FormValue("state")
	redirectUrl := session.GetString(state)
	if redirectUrl == "" {
		ctx.WriteString("invalid state")
		ctx.StatusCode(http.StatusBadRequest)
		return
	}
	session.Delete(state)
	ctx.RemoveCookie(x.Options.Coki_Session)
	session.Destroy()

	// 跳转回登出时的页面
	ctx.Redirect(redirectUrl, http.StatusFound)
}

func (x *defaultOIDCClient) SaveToken(ctx context.Context, token *oauth2.Token) error {
	j, err := sjson.Serialize(token)
	if u.LogError(err) {
		return err
	}

	// 保存常规令牌
	session := x.Options.SessionManager.Start(ctx)
	session.Set(COKI_TOKEN, j)

	// 保存ID令牌
	idToken, ok := token.Extra("id_token").(string)
	if ok {
		session.Set(COKI_IDTOKEN, idToken)
	}
	// err = x.Options.SecureCookie.Set(ctx, COKI_TOKEN, j)
	// if u.LogError(err) {
	// 	return err
	// }

	// // 保存ID令牌
	// if idToken != "" {
	// 	err = x.Options.SecureCookie.Set(ctx, COKI_IDTOKEN, idToken)
	// 	return err
	// }

	return nil
}

func (x *defaultOIDCClient) GetToken(ctx context.Context) (*oauth2.Token, string, error) {
	session := x.Options.SessionManager.Start(ctx)
	j := session.GetString(COKI_TOKEN)
	// j, err := x.Options.SecureCookie.Get(ctx, COKI_TOKEN)
	// if u.LogError(err) {
	// 	return nil, err
	// }

	t := new(oauth2.Token)
	err := sjson.Deserialize(j, t)
	if u.LogError(err) {
		return nil, "", err
	}

	idToken := session.GetString(COKI_IDTOKEN)

	return t, idToken, err
}

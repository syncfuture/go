package soidc

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/syncfuture/go/json"
	u "github.com/syncfuture/go/util"

	"github.com/syncfuture/go/rand"

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
}

func NewOIDCClient(options *ClientOptions) IOIDCClient {
	checkOptions(options)

	x := new(defaultOIDCClient)
	x.Options = options

	ctx := gocontext.Background()
	var err error
	x.OIDCProvider, err = oidc.NewProvider(ctx, options.ProviderUrl)
	if err != nil {
		log.Fatal(err)
	}
	x.OAuth2Config = new(oauth2.Config)
	x.OAuth2Config.ClientID = options.ClientID
	x.OAuth2Config.ClientSecret = options.ClientSecret
	x.OAuth2Config.Endpoint = x.OIDCProvider.Endpoint()
	x.OAuth2Config.RedirectURL = options.CallbackURL
	x.OAuth2Config.Scopes = append(options.Scopes, oidc.ScopeOpenID)

	return x
}

func (x *defaultOIDCClient) HandleAuthentication(ctx context.Context) {
	session := x.Options.Sessions.Start(ctx)

	handlerName := ctx.GetCurrentRoute().MainHandlerName()
	area, controller, action := getRoutes(handlerName)

	// 判断请求是否允许访问
	userid := session.GetString(x.Options.Sess_ID)
	if userid != "" {
		roles := session.GetInt64Default(x.Options.Sess_Roles, 0)
		// 已登录
		allow := x.Options.PermissionAuditor.CheckRoute(area, controller, action, roles)
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
		allow := x.Options.PermissionAuditor.CheckRoute(area, controller, action, 0)
		if allow {
			// 允许匿名
			ctx.Next()
			return
		}

		// 记录请求地址，跳转去登录页面
		state := rand.String(32)
		session.Set(state, ctx.Request().URL.String())
		ctx.Redirect(x.OAuth2Config.AuthCodeURL(state), http.StatusFound)
	}
}

func (x *defaultOIDCClient) getRoutes(handlerName string) (string, string, string) {
	array := strings.Split(handlerName, ".")
	return array[0], array[1], array[2]
}

func (x *defaultOIDCClient) HandleSignInCallback(ctx context.Context) {
	session := x.Options.Sessions.Start(ctx)

	state := ctx.FormValue("state")
	redirectUrl := session.GetString(state)
	if redirectUrl == "" {
		ctx.WriteString("invalid state")
		ctx.StatusCode(http.StatusBadRequest)
		return
	}
	session.Delete(x.Options.Sess_State) // 释放内存

	// 交换令牌
	code := ctx.FormValue("code")
	httpCtx := gocontext.Background()
	oauth2Token, err := x.OAuth2Config.Exchange(httpCtx, code /*,oauth2.SetAuthURLParam("aaa","bbb")*/)
	if err != nil {
		ctx.Write([]byte("Failed to exchange token: " + err.Error()))
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	userInfo, err := x.OIDCProvider.UserInfo(httpCtx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		ctx.Write([]byte("Failed to get userinfo: " + err.Error()))
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	claims := make(map[string]string, 6)
	userInfo.Claims(&claims)

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
}

func (x *defaultOIDCClient) HandleSignOutCallback(ctx context.Context) {
	session := x.Options.Sessions.Start(ctx)

	session.Destroy()
	ctx.RemoveCookie(x.Options.Coki_Token)
	ctx.RemoveCookie(x.Options.Coki_Session)

	// Todo: 去Passport注销
	ctx.Redirect("/", http.StatusFound)
}

func getRoutes(handlerName string) (string, string, string) {
	array := strings.Split(handlerName, ".")
	return array[0], array[1], array[2]
}

func checkOptions(options *ClientOptions) {
	if options.ClientID == "" {
		log.Fatal("OIDCClient.Options.ClientID cannot be empty.")
	}
	if options.ClientSecret == "" {
		log.Fatal("OIDCClient.Options.ClientSecret cannot be empty.")
	}
	if len(options.Scopes) == 0 {
		log.Fatal("OIDCClient.Options.Scopes cannot be empty")
	}
	if options.ProviderUrl == "" {
		log.Fatal("OIDCClient.Options.ProviderUrl cannot be empty.")
	}
	if options.CallbackURL == "" {
		log.Fatal("OIDCClient.Options.CallbackURL cannot be empty.")
	}

	if options.Sessions == nil {
		log.Fatal("OIDCClient.Options.Sessions cannot be nil")
	}
	if options.PermissionAuditor == nil {
		log.Fatal("OIDCClient.Options.PermissionAuditor cannot be nil")
	}
	if options.SecureCookie == nil {
		log.Fatal("OIDCClient.Options.SecureCookie cannot be nil")
	}

	if options.Coki_Token == "" {
		options.Coki_Token = COKI_TOKEN
	}
	if options.Coki_Session == "" {
		options.Coki_Session = COKI_SESSION
	}

	if options.Sess_ID == "" {
		options.Sess_ID = SESS_ID
	}
	if options.Sess_Username == "" {
		options.Sess_Username = SESS_USERNAME
	}
	if options.Sess_Email == "" {
		options.Sess_Email = SESS_EMAIL
	}
	if options.Sess_Roles == "" {
		options.Sess_Roles = SESS_ROLES
	}
	if options.Sess_Level == "" {
		options.Sess_Level = SESS_LEVEL
	}
	if options.Sess_Status == "" {
		options.Sess_Status = SESS_STATUS
	}

	if options.Sess_State == "" {
		options.Sess_State = SESS_STATE
	}

	if options.AccessDeniedURL == "" {
		options.AccessDeniedURL = "/"
	}
}

func (x *defaultOIDCClient) NewHttpClient(ctx context.Context) (*http.Client, error) {
	session := x.Options.Sessions.Start(ctx)
	userID := session.GetString(SESS_ID)
	if userID == "" {
		return nil, fmt.Errorf("user id not exists in session")
	}

	t, err := x.GetToken(ctx)
	if u.LogError(err) {
		return nil, err
	}

	goctx := gocontext.Background()
	tokenSource := x.OAuth2Config.TokenSource(goctx, t)
	newToken, err := tokenSource.Token()
	if u.LogError(err) {
		return nil, err
	}

	if newToken.AccessToken != t.AccessToken {
		x.SaveToken(ctx, newToken)
		log.Debugf("Saved new token for user %s", userID)
	}

	return oauth2.NewClient(goctx, tokenSource), nil
}

func (x *defaultOIDCClient) SaveToken(ctx context.Context, token *oauth2.Token) error {
	j, err := json.Serialize(token)
	if u.LogError(err) {
		return err
	}

	err = x.Options.SecureCookie.Set(ctx, COKI_TOKEN, j)
	return err
}

func (x *defaultOIDCClient) GetToken(ctx context.Context) (*oauth2.Token, error) {
	j, err := x.Options.SecureCookie.Get(ctx, COKI_TOKEN)
	if u.LogError(err) {
		return nil, err
	}

	t := new(oauth2.Token)
	err = json.Deserialize(j, t)
	u.LogError(err)
	return t, err
}
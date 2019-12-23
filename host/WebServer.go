package host

import (
	"time"

	log "github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/view"
	"github.com/syncfuture/go/config"
	"github.com/syncfuture/go/security"
	"github.com/syncfuture/go/soidc"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/go/surl"
)

type WebServer struct {
	ConfigProvider          config.IConfigProvider
	App                     *iris.Application
	SessionManager          *sessions.Sessions
	URLProvider             surl.IURLProvider
	RoutePermissionProvider security.IRoutePermissionProvider
	PermissionAuditor       security.IPermissionAuditor
	OIDCClient              soidc.IOIDCClient
	RedisConfig             *sredis.RedisConfig
	SecureCookie            security.ISecureCookie
	ViewEngine              view.Engine
	StaticFilesDir          string
	IsDebug                 bool
}

func NewWebServer() (r *WebServer) {
	r = new(WebServer)
	r.ConfigProvider = config.NewJsonConfigProvider()
	// 日志和配置
	logLevel := r.ConfigProvider.GetStringDefault("Log.Level", "warn")
	log.SetLevel(logLevel)

	// Http客户端
	ConfigHttpClient(r)

	// Redis
	r.RedisConfig = r.ConfigProvider.GetRedisConfig()

	// URLProvider
	r.URLProvider = surl.NewRedisURLProvider(r.RedisConfig)

	// 权限
	projectName := r.ConfigProvider.GetString("ProjectName")
	if projectName == "" {
		log.Fatal("cannot find 'ProjectName' config")
	}
	r.RoutePermissionProvider = security.NewRedisRoutePermissionProvider(projectName, r.RedisConfig)
	r.PermissionAuditor = security.NewPermissionAuditor(r.RoutePermissionProvider)

	// 渲染URL
	oidcConfig := r.ConfigProvider.GetOIDCConfig()
	oidcConfig.PassportURL = r.URLProvider.RenderURLCache(oidcConfig.PassportURL)
	oidcConfig.SignInCallbackURL = r.URLProvider.RenderURLCache(oidcConfig.SignInCallbackURL)
	oidcConfig.SignOutCallbackURL = r.URLProvider.RenderURLCache(oidcConfig.SignOutCallbackURL)

	// Cookie和Session
	hashKey := r.ConfigProvider.GetString("Cookie.HashKey")
	if hashKey == "" {
		log.Fatal("Cannot find 'Cookie.HashKey' config")
	}
	blockKey := r.ConfigProvider.GetString("Cookie.BlockKey")
	if blockKey == "" {
		log.Fatal("Cannot find 'Cookie.BlockKey' config")
	}
	r.SecureCookie = security.NewDefaultSecureCookie([]byte(hashKey), []byte(blockKey))
	r.SessionManager = sessions.New(sessions.Config{
		Expires:      -1 * time.Hour,
		Cookie:       soidc.COKI_SESSION,
		Encode:       r.SecureCookie.Encode,
		Decode:       r.SecureCookie.Decode,
		AllowReclaim: true,
	})

	// OIDC
	oidcOptions := &soidc.ClientOptions{
		ClientID:           oidcConfig.ClientID,
		ClientSecret:       oidcConfig.ClientSecret,
		PassportURL:        oidcConfig.PassportURL,
		SignInCallbackURL:  oidcConfig.SignInCallbackURL,
		SignOutCallbackURL: oidcConfig.SignOutCallbackURL,
		AccessDeniedURL:    oidcConfig.AccessDeniedURL,
		Scopes:             oidcConfig.Scopes,
		Coki_Session:       soidc.COKI_SESSION,
		PermissionAuditor:  r.PermissionAuditor,
		SessionManager:     r.SessionManager,
	}
	r.OIDCClient = soidc.NewOIDCClient(oidcOptions)

	// IRIS App
	r.App = iris.New()
	r.App.Logger().SetLevel(logLevel)
	r.App.Use(recover.New())
	r.App.Use(logger.New())

	// 登录登出相关
	r.App.Get("/signin", r.OIDCClient.HandleSignIn)
	r.App.Get("/signin-oidc", r.OIDCClient.HandleSignInCallback)
	r.App.Get("/signout", r.OIDCClient.HandleSignOut)
	r.App.Get("/signout-callback-oidc", r.OIDCClient.HandleSignOutCallback)

	// Debug模式
	r.IsDebug = r.ConfigProvider.GetBool("Dev.Debug")
	// 视图引擎
	if r.ViewEngine == nil {
		r.ViewEngine = iris.HTML("./views", ".html").Layout("shared/_layout.html").Reload(r.IsDebug)
	}
	r.App.RegisterView(r.ViewEngine)

	// 网站静态文件根目录
	if r.StaticFilesDir == "" {
		r.StaticFilesDir = "./wwwroot"
	}
	r.App.HandleDir("/", r.StaticFilesDir)

	return r
}

func (x *WebServer) GetConfigProvider() config.IConfigProvider {
	return x.ConfigProvider
}

func (x *WebServer) Run() {
	listenAddr := x.ConfigProvider.GetString("ListenAddr")
	if listenAddr == "" {
		log.Fatal("Cannot find 'ListenAddr' config")
	}
	x.App.Run(iris.Addr(listenAddr))
}

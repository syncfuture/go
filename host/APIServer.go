package host

import (
	"crypto/tls"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/jwt"
	log "github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/syncfuture/go/config"
	"github.com/syncfuture/go/security"
	"github.com/syncfuture/go/soidc"
	"github.com/syncfuture/go/surl"
	"net/http"
	"net/url"
	"strings"
)

type APIServerOption struct {
	ActionMap *map[string]*Action
}

type APIServer struct {
	ConfigProvider          config.IConfigProvider
	App                     *iris.Application
	PublicKeyProvider       soidc.IPublicKeyProvider
	URLProvider             surl.IURLProvider
	RoutePermissionProvider security.IRoutePermissionProvider
	PermissionAuditor       security.IPermissionAuditor
	PreMiddlewares          []context.Handler
	ActionMap               *map[string]*Action
}

func NewApiServer(option *APIServerOption) (r *APIServer) {
	r.ConfigProvider = config.NewJsonConfigProvider()
	// 日志和配置
	logLevel := r.ConfigProvider.GetString("Log.Level")
	log.SetLevel(logLevel)

	// 调试
	isDebug := r.ConfigProvider.GetBool("Dev.Debug")
	if isDebug {
		// 跳过证书验证
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		// 使用代理
		proxy := r.ConfigProvider.GetString("Dev.Proxy")
		if proxy != "" {
			transport.Proxy = func(r *http.Request) (*url.URL, error) {
				return url.Parse(proxy)
			}
		}
		http.DefaultClient.Transport = transport
	}

	// URLProvider
	redisConfig := r.ConfigProvider.GetRedisConfig()
	r.URLProvider = surl.NewRedisURLProvider(redisConfig)

	oidcConfig := r.ConfigProvider.GetOIDCConfig()
	jwksURL := oidcConfig.JWKSURL
	if jwksURL == "" {
		jwksURL = "/.well-known/openid-configuration/jwks"
	}

	// 权限
	projectName := r.ConfigProvider.GetString("ProjectName")
	r.RoutePermissionProvider = security.NewRedisRoutePermissionProvider(projectName, redisConfig)
	r.PermissionAuditor = security.NewPermissionAuditor(r.RoutePermissionProvider)

	// 渲染URL
	oidcConfig.PassportURL = r.URLProvider.RenderURLCache(oidcConfig.PassportURL)

	// 公钥提供器
	r.PublicKeyProvider = soidc.NewPublicKeyProvider(oidcConfig.PassportURL, jwksURL, projectName)

	// JWT验证中间件
	jwtMiddleware := jwt.New(jwt.Config{
		ValidationKeyGetter: r.PublicKeyProvider.GetKey,
		SigningMethod:       jwtgo.SigningMethodRS256,
	})

	// 授权中间件
	authMiddleware := &AuthMidleware{
		ActionMap:         r.ActionMap,
		PermissionAuditor: r.PermissionAuditor,
	}

	// 添加中间件
	r.PreMiddlewares = append(r.PreMiddlewares, jwtMiddleware.Serve)
	r.PreMiddlewares = append(r.PreMiddlewares, authMiddleware.Serve)

	// IRIS App
	r.App = iris.New()
	r.App.Logger().SetLevel(logLevel)
	r.App.Use(recover.New())
	r.App.Use(logger.New())

	r.ActionMap = option.ActionMap

	return r
}

func (x *APIServer) Run() {
	x.registerActions()

	listenAddr := x.ConfigProvider.GetString("ListenAddr")
	x.App.Run(iris.Addr(listenAddr))
}

func (x *APIServer) registerActions() {
	for name, action := range *x.ActionMap {
		handlers := append(x.PreMiddlewares, action.Handler)
		x.registerAction(name, handlers...)
	}
}

func (x *APIServer) registerAction(name string, handlers ...context.Handler) {
	index := strings.Index(name, "/")
	method := name[:index]
	path := name[index:]

	switch method {
	case http.MethodPost:
		x.App.Post(path, handlers...)
		break
	case http.MethodGet:
		x.App.Get(path, handlers...)
		break
	case http.MethodPut:
		x.App.Put(path, handlers...)
		break
	case http.MethodPatch:
		x.App.Patch(path, handlers...)
		break
	case http.MethodDelete:
		x.App.Delete(path, handlers...)
		break
	default:
		panic("does not support method " + method)
	}
}

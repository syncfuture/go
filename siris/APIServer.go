package siris

import (
	"net/http"
	"strings"

	"github.com/syncfuture/go/security"

	"github.com/syncfuture/go/sredis"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/syncfuture/go/soidc"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type (
	APIServer struct {
		listenAddr     string
		App            *iris.Application
		preMiddlewares []context.Handler
		actionMap      *map[string]*Action
	}
)

func NewAPIServer(
	projectName, logLevel, listenAddr string,
	redisConfig *sredis.RedisConfig,
	soidcConfig *soidc.OIDCConfig,
	actionMap *map[string]*Action,
) *APIServer {
	r := new(APIServer)

	r.App = iris.New()
	r.App.Logger().SetLevel(logLevel)

	r.App.Use(recover.New())
	r.App.Use(logger.New())

	r.actionMap = actionMap
	r.listenAddr = listenAddr

	jwksURL := soidcConfig.JWKSURL
	if soidcConfig.JWKSURL == "" {
		jwksURL = "/.well-known/openid-configuration/jwks"
	}
	publicKeyProvider := soidc.NewPublicKeyProvider(soidcConfig.PassportURL, jwksURL, projectName)
	routePermissionProvider := security.NewRedisRoutePermissionProvider(projectName, redisConfig)
	permissionAuditor := security.NewPermissionAuditor(routePermissionProvider)

	jwtMiddleware := jwt.New(jwt.Config{
		ValidationKeyGetter: publicKeyProvider.GetKey,
		SigningMethod:       jwtgo.SigningMethodRS256,
	})

	authMiddleware := &AuthMidleware{
		ActionMap:         r.actionMap,
		PermissionAuditor: permissionAuditor,
	}

	r.preMiddlewares = append(r.preMiddlewares, jwtMiddleware.Serve)
	r.preMiddlewares = append(r.preMiddlewares, authMiddleware.Serve)

	return r
}

func (x *APIServer) Run() {
	x.registerActions()
	x.App.Run(iris.Addr(x.listenAddr))
}

func (x *APIServer) registerActions() {
	for name, action := range *x.actionMap {
		handlers := append(x.preMiddlewares, action.Handler)
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
	case http.MethodDelete:
		x.App.Delete(path, handlers...)
		break
	default:
		panic("does not support method " + method)
	}
}

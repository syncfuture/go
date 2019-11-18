package siris

import (
	"github.com/syncfuture/go/soidc"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type (
	WebServer struct {
		listenAddr     string
		App            *iris.Application
		preMiddlewares []context.Handler
		actionMap      *map[string]*Action
	}
)

func NewWebServer(
	logLevel, listenAddr string,
	debug bool,
	oidcClient soidc.IOIDCClient,
) *WebServer {
	r := new(WebServer)

	r.app = iris.New()
	r.app.Logger().SetLevel(logLevel)

	r.app.Use(recover.New())
	r.app.Use(logger.New())

	r.listenAddr = listenAddr

	r.app.Get("/signin-oidc", oidcClient.HandleSignInCallback)
	r.app.Get("/signout", oidcClient.HandleSignOut)
	r.app.Get("/signout-callback-oidc", oidcClient.HandleSignOutCallback)

	viewEngine := iris.HTML("./views", ".html").Layout("shared/_layout.html").Reload(debug)
	r.app.RegisterView(viewEngine)
	r.app.HandleDir("/", "./wwwroot")
	r.app.Use(oidcClient.HandleAuthentication)

	return r
}

func (x *WebServer) Run() {
	x.app.Run(iris.Addr(x.listenAddr))
}

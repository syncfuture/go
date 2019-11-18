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

	r.App = iris.New()
	r.App.Logger().SetLevel(logLevel)

	r.App.Use(recover.New())
	r.App.Use(logger.New())

	r.listenAddr = listenAddr

	r.App.Get("/signin-oidc", oidcClient.HandleSignInCallback)
	r.App.Get("/signout", oidcClient.HandleSignOut)
	r.App.Get("/signout-callback-oidc", oidcClient.HandleSignOutCallback)

	viewEngine := iris.HTML("./views", ".html").Layout("shared/_layout.html").Reload(debug)
	r.App.RegisterView(viewEngine)
	r.App.HandleDir("/", "./wwwroot")
	r.App.Use(oidcClient.HandleAuthentication)

	return r
}

func (x *WebServer) Run() {
	x.App.Run(iris.Addr(x.listenAddr))
}

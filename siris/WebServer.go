package siris

import (
	"github.com/syncfuture/go/soidc"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type (
	WebServerOption struct {
		ListenAddr, LogLevel string
		Debug                bool
		OIDCClient           soidc.IOIDCClient
	}
	WebServer struct {
		listenAddr     string
		App            *iris.Application
		preMiddlewares []context.Handler
		actionMap      *map[string]*Action
	}
)

func NewWebServer(option *WebServerOption) *WebServer {
	r := new(WebServer)

	r.App = iris.New()
	r.App.Logger().SetLevel(option.LogLevel)

	r.App.Use(recover.New())
	r.App.Use(logger.New())

	r.listenAddr = option.ListenAddr

	r.App.Get("/signin-oidc", option.OIDCClient.HandleSignInCallback)
	r.App.Get("/signout", option.OIDCClient.HandleSignOut)
	r.App.Get("/signout-callback-oidc", option.OIDCClient.HandleSignOutCallback)

	viewEngine := iris.HTML("./views", ".html").Layout("shared/_layout.html").Reload(option.Debug)
	r.App.RegisterView(viewEngine)
	r.App.HandleDir("/", "./wwwroot")

	return r
}

func (x *WebServer) Run() {
	x.App.Run(iris.Addr(x.listenAddr))
}

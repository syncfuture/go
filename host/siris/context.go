package siris

import (
	"net/url"
	"time"

	"github.com/Lukiya/home/go/web/core"
	"github.com/kataras/iris/v12"
	"github.com/syncfuture/go/soidc"
)

func InitContext(ctx iris.Context) {
	session := core.WebServer.SessionManager.Start(ctx)
	ctx.ViewData("Year", time.Now().Year())
	ctx.ViewData("Debug", core.WebServer.ConfigProvider.GetBool("Dev.Debug"))
	ctx.ViewData("Version", core.WebServer.ConfigProvider.GetStringDefault("Version", "1.0.0"))

	ctx.ViewData("CurrentURL", url.PathEscape(ctx.Request().URL.String()))
	userID := session.GetString(soidc.SESS_ID)
	ctx.ViewData("UserID", userID)
	isAuth := userID != ""
	ctx.ViewData("IsAuthenticated", isAuth)
	if isAuth {
		ctx.ViewData("Username", session.GetString(soidc.SESS_USERNAME))
		ctx.ViewData("UserRoles", session.GetInt64Default(soidc.SESS_ROLES, 0))
	}
}

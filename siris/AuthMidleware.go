package siris

import (
	"net/http"
	"strconv"

	"github.com/syncfuture/go/security"

	u "github.com/syncfuture/go/util"

	log "github.com/kataras/golog"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12/context"
)

type AuthMidleware struct {
	PermissionAuditor security.IPermissionAuditor
	ActionMap         *map[string]*Action
}

func (x *AuthMidleware) Serve(ctx context.Context) {
	var msgCode string
	token := ctx.Values().Get("jwt").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	if roleStr, ok := claims["role"].(string); ok && roleStr != "" {
		// 有角色
		roles, err := strconv.ParseInt(roleStr, 10, 64)
		if !u.LogError(err) {
			// 角色是数字
			route := ctx.GetCurrentRoute().Name()
			if action, ok := (*x.ActionMap)[route]; ok {
				// 找到了对应的Action
				if x.PermissionAuditor.CheckRoute(action.Area, action.Controller, action.Action, roles) {
					// 有权限，允许访问
					ctx.Next()
					return
				} else {
					msgCode = "permission denied"
				}
			} else {
				msgCode = route + " doesn't exist in action map"
				log.Warn(msgCode)
			}
		} else {
			msgCode = "parse role error"
		}
	} else {
		msgCode = "token doesn't have role field"
		log.Warn(msgCode, " ", claims)
	}

	// 无权限
	ctx.StatusCode(http.StatusUnauthorized)
	ctx.WriteString(msgCode)
}

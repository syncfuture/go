package host

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/kataras/iris/v12/context"
	u "github.com/syncfuture/go/util"
)

type (
	Action struct {
		Route      string
		Area       string
		Controller string
		Action     string
		Handler    context.Handler
	}
)

func NewAction(route, area, controller string, handler context.Handler) *Action {
	action := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	action = action[strings.LastIndex(action, ".")+1:]

	return &Action{
		Route:      route,
		Area:       area,
		Controller: controller,
		Action:     action,
		Handler:    handler,
	}
}

func CreateActionMap(actionGroups ...*[]*Action) *map[string]*Action {
	actionMap := make(map[string]*Action)

	for _, actionGroup := range actionGroups {
		for _, action := range *actionGroup {
			actionMap[action.Route] = action
		}
	}

	return &actionMap
}

func HandleError(ctx context.Context, err error) bool {
	if u.LogError(err) {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return true
	}
	return false
}

package host

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"

	"github.com/syncfuture/go/config"

	"github.com/kataras/iris/v12/context"
	u "github.com/syncfuture/go/util"
)

type IServer interface {
	GetConfigProvider() config.IConfigProvider
}

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

func HandleError(ctx context.Context, err error) bool {
	if u.LogError(err) {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return true
	}
	return false
}

func ConfigHttpClient(x IServer) {
	configProvider := x.GetConfigProvider()
	// Http客户端配置
	skipCertVerification := configProvider.GetBool("Http.SkipCertVerification")
	proxy := configProvider.GetString("Http.Proxy")
	if skipCertVerification || proxy != "" {
		// 任意条件满足，则使用自定义传输层
		transport := new(http.Transport)
		if skipCertVerification {
			// 跳过证书验证
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: skipCertVerification}
		}
		if proxy != "" {
			// 使用代理
			transport.Proxy = func(r *http.Request) (*url.URL, error) {
				return url.Parse(proxy)
			}
		}
		http.DefaultClient.Transport = transport
	}
}

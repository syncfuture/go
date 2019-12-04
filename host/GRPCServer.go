package host

import (
	"crypto/tls"
	log "github.com/kataras/golog"
	panichandler "github.com/kazegusuri/grpc-panic-handler"
	"github.com/syncfuture/go/config"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"net/url"
)

type GRPCServer struct {
	ConfigProvider config.IConfigProvider
	Server         *grpc.Server
}

func NewGRPCServer() (r *GRPCServer) {
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

	// GRPC Server
	uIntOpt := grpc.UnaryInterceptor(panichandler.UnaryPanicHandler)
	sIntOpt := grpc.StreamInterceptor(panichandler.StreamPanicHandler)
	panichandler.InstallPanicHandler(func(r interface{}) {
		log.Error(r)
	})
	r.Server = grpc.NewServer(uIntOpt, sIntOpt)

	return r
}

func (x *GRPCServer) Run() {
	listenAddr := x.ConfigProvider.GetString("ListenAddr")
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := x.Server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

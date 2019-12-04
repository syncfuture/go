package host

import (
	"crypto/tls"
	log "github.com/kataras/golog"
	panichandler "github.com/kazegusuri/grpc-panic-handler"
	"github.com/syncfuture/go/config"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/go/surl"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"net/url"
)

type ServiceServer struct {
	ConfigProvider config.IConfigProvider
	GRPCServer     *grpc.Server
	URLProvider    surl.IURLProvider
	RedisConfig    *sredis.RedisConfig
}

func NewGRPCServer() (r *ServiceServer) {
	r = new(ServiceServer)
	r.ConfigProvider = config.NewJsonConfigProvider()
	// 日志和配置
	logLevel := r.ConfigProvider.GetStringDefault("Log.Level", "warn")
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

	// Redis
	r.RedisConfig = r.ConfigProvider.GetRedisConfig()

	// URLProvider
	r.URLProvider = surl.NewRedisURLProvider(r.RedisConfig)

	// GRPC Server
	uIntOpt := grpc.UnaryInterceptor(panichandler.UnaryPanicHandler)
	sIntOpt := grpc.StreamInterceptor(panichandler.StreamPanicHandler)
	panichandler.InstallPanicHandler(func(r interface{}) {
		log.Error(r)
	})
	r.GRPCServer = grpc.NewServer(uIntOpt, sIntOpt)

	return r
}

func (x *ServiceServer) Run() {
	listenAddr := x.ConfigProvider.GetString("ListenAddr")
	if listenAddr == "" {
		log.Fatal("Cannot find 'ListenAddr' config")
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Infof("Listening at %v\n", listenAddr)
	if err := x.GRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package host

import (
	"net"

	"github.com/syncfuture/go/sgrpc"

	log "github.com/kataras/golog"
	panichandler "github.com/kazegusuri/grpc-panic-handler"
	"github.com/syncfuture/go/config"
	"github.com/syncfuture/go/security"
	"github.com/syncfuture/go/sredis"
	"google.golang.org/grpc"
)

type ServiceServer struct {
	ConfigProvider    config.IConfigProvider
	RedisConfig       *sredis.RedisConfig
	PermissionAuditor security.IPermissionAuditor
	GRPCServer        *grpc.Server
}

func NewServiceServer() (r *ServiceServer) {
	r = new(ServiceServer)
	r.ConfigProvider = config.NewJsonConfigProvider()
	// 日志和配置
	logLevel := r.ConfigProvider.GetStringDefault("Log.Level", "warn")
	log.SetLevel(logLevel)

	// Http客户端
	ConfigHttpClient(r)

	// Redis
	r.RedisConfig = r.ConfigProvider.GetRedisConfig()

	// 权限
	if len(r.RedisConfig.Addrs) > 0 {
		routePermissionProvider := security.NewRedisRoutePermissionProvider("", r.RedisConfig)
		r.PermissionAuditor = security.NewPermissionAuditor(routePermissionProvider)
	}

	// GRPC Server
	unaryPanicHandler := grpc.UnaryInterceptor(panichandler.UnaryPanicHandler)
	streamPanicHandler := grpc.StreamInterceptor(panichandler.StreamPanicHandler)
	panichandler.InstallPanicHandler(func(r interface{}) {
		log.Error(r)
	})
	unaryJWTHandler := grpc.UnaryInterceptor(sgrpc.AttachJWTToken) // 负责附加jwt令牌到上下文
	r.GRPCServer = grpc.NewServer(unaryPanicHandler, streamPanicHandler, unaryJWTHandler)

	return r
}

func (x *ServiceServer) GetConfigProvider() config.IConfigProvider {
	return x.ConfigProvider
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

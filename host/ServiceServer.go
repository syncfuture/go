package host

import (
	"net"

	"github.com/syncfuture/go/sgrpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	panichandler "github.com/kazegusuri/grpc-panic-handler"
	"github.com/syncfuture/go/config"
	"github.com/syncfuture/go/security"
	log "github.com/syncfuture/go/slog"
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
	log.Init()

	// Http客户端
	ConfigHttpClient(r)

	// Redis
	r.ConfigProvider.GetStruct("Redis", &r.RedisConfig)

	// 权限
	if r.RedisConfig != nil && len(r.RedisConfig.Addrs) > 0 {
		routePermissionProvider := security.NewRedisRoutePermissionProvider("", r.RedisConfig)
		r.PermissionAuditor = security.NewPermissionAuditor(routePermissionProvider)
	}

	// GRPC Server
	unaryHandler := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(panichandler.UnaryPanicHandler, sgrpc.AttachJWTToken))
	streamHandler := grpc.StreamInterceptor(panichandler.StreamPanicHandler)
	panichandler.InstallPanicHandler(func(r interface{}) {
		log.Error(r)
	})
	maxSize := r.ConfigProvider.GetIntDefault("GRPC.MaxRecvMsgSize", 10*1024*1024)
	r.GRPCServer = grpc.NewServer(grpc.MaxRecvMsgSize(maxSize), unaryHandler, streamHandler)

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

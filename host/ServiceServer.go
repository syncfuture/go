package host

import (
	"net"

	log "github.com/kataras/golog"
	panichandler "github.com/kazegusuri/grpc-panic-handler"
	"github.com/syncfuture/go/config"
	"google.golang.org/grpc"
)

type ServiceServer struct {
	ConfigProvider config.IConfigProvider
	GRPCServer     *grpc.Server
}

func NewServiceServer() (r *ServiceServer) {
	r = new(ServiceServer)
	r.ConfigProvider = config.NewJsonConfigProvider()
	// 日志和配置
	logLevel := r.ConfigProvider.GetStringDefault("Log.Level", "warn")
	log.SetLevel(logLevel)

	// Http客户端
	ConfigHttpClient(r)

	// GRPC Server
	uIntOpt := grpc.UnaryInterceptor(panichandler.UnaryPanicHandler)
	sIntOpt := grpc.StreamInterceptor(panichandler.StreamPanicHandler)
	panichandler.InstallPanicHandler(func(r interface{}) {
		log.Error(r)
	})
	r.GRPCServer = grpc.NewServer(uIntOpt, sIntOpt)

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

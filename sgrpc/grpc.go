package sgrpc

import (
	log "github.com/kataras/golog"
	panichandler "github.com/kazegusuri/grpc-panic-handler"
	"google.golang.org/grpc"
)

func CreateServer() *grpc.Server {
	uIntOpt := grpc.UnaryInterceptor(panichandler.UnaryPanicHandler)
	sIntOpt := grpc.StreamInterceptor(panichandler.StreamPanicHandler)
	panichandler.InstallPanicHandler(func(r interface{}) {
		log.Error(r)
	})
	return grpc.NewServer(uIntOpt, sIntOpt)
}

package sgrpc

import (
	panichandler "github.com/kazegusuri/grpc-panic-handler"
	log "github.com/syncfuture/go/slog"
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

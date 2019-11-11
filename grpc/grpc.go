package grpc

import (
	log "github.com/kataras/golog"
	panichandler "github.com/kazegusuri/grpc-panic-handler"
	g "google.golang.org/grpc"
)

func CreateServer() *g.Server {
	uIntOpt := g.UnaryInterceptor(panichandler.UnaryPanicHandler)
	sIntOpt := g.StreamInterceptor(panichandler.StreamPanicHandler)
	panichandler.InstallPanicHandler(func(r interface{}) {
		log.Error(r)
	})
	return g.NewServer(uIntOpt, sIntOpt)
}

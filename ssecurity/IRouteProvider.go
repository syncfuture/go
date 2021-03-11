package ssecurity

import "github.com/syncfuture/go/sproto"

type IRouteProvider interface {
	CreateRoute(*sproto.RouteDTO) error
	GetRoute(string) (*sproto.RouteDTO, error)
	UpdateRoute(*sproto.RouteDTO) error
	RemoveRoute(string) error
	GetRoutes() (map[string]*sproto.RouteDTO, error)
}

package security

import (
	"github.com/syncfuture/go/sproto"
)

type IRoutePermissionProvider interface {
	CreateRoute(*sproto.RouteDTO) error
	GetRoute(string) (*sproto.RouteDTO, error)
	UpdateRoute(*sproto.RouteDTO) error
	RemoveRoute(string) error
	GetRoutes() (map[string]*sproto.RouteDTO, error)

	CreatePermission(*sproto.PermissionDTO) error
	GetPermission(string) (*sproto.PermissionDTO, error)
	UpdatePermission(*sproto.PermissionDTO) error
	RemovePermission(string) error
	GetPermissions() (map[string]*sproto.PermissionDTO, error)
}

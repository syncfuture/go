package ssecurity

import "github.com/syncfuture/go/sproto"

type IPermissionProvider interface {
	CreatePermission(*sproto.PermissionDTO) error
	GetPermission(string) (*sproto.PermissionDTO, error)
	UpdatePermission(*sproto.PermissionDTO) error
	RemovePermission(string) error
	GetPermissions() (map[string]*sproto.PermissionDTO, error)
}

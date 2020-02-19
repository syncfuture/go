package security

import (
	log "github.com/kataras/golog"
	"github.com/syncfuture/go/sproto"
)

type IPermissionAuditor interface {
	CheckPermission(permissionID string, userRoles int64, userLevel int32) bool
	CheckRoute(area, controller, action string, userRoles int64, userLevel int32) bool
}

type permissionAuditor struct {
	routePermissionProvider IRoutePermissionProvider
	routes                  map[string]*sproto.RouteDTO
	permissions             map[string]*sproto.PermissionDTO
}

func NewPermissionAuditor(routePermissionProvider IRoutePermissionProvider) IPermissionAuditor {
	r := new(permissionAuditor)
	r.routePermissionProvider = routePermissionProvider
	r.ReloadRoutePermissions()
	return r
}

func (x *permissionAuditor) ReloadRoutePermissions() error {
	var err error

	x.routes, err = x.routePermissionProvider.GetRoutes()
	if err != nil {
		return err
	}
	x.permissions, err = x.routePermissionProvider.GetPermissions()
	if err != nil {
		return err
	}

	return nil
}

func (x *permissionAuditor) CheckPermission(permissionID string, userRoles int64, userLevel int32) bool {
	if permission, exists := x.permissions[permissionID]; exists {
		return checkPermission(permission, userRoles, userLevel)
	}

	log.Warnf("permission: %s does not exist", permissionID)
	return false
}

func checkPermission(permission *sproto.PermissionDTO, userRoles int64, userLevel int32) bool {
	if permission.IsAllowGuest {
		return true
	} else if permission.IsAllowAnyUser {
		return userRoles > 0
	} else {
		return (permission.AllowedRoles&userRoles) > 0 && userLevel > permission.Level
	}
}

func (x *permissionAuditor) CheckRoute(area, controller, action string, userRoles int64, userLevel int32) bool {
	key := area + "_" + controller + "_" + action

	route := new(sproto.RouteDTO)
	var exists bool
	if route, exists = x.routes[key]; !exists {
		key = area + "_" + controller + "_"
		if route, exists = x.routes[key]; !exists {
			key = area + "__"
			if route, exists = x.routes[key]; !exists {
				log.Warnf("route: [%s,%s,%s] does not exist", area, controller, action)
				return false
			}
		}
	}

	if permission, exists := x.permissions[route.Permission_ID]; exists {
		return checkPermission(permission, userRoles, userLevel)
	}

	log.Warnf("permission: %s does not exist", route.Permission_ID)
	return false
}

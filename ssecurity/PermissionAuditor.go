package ssecurity

import (
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/sproto"
	"github.com/syncfuture/go/u"
)

type IPermissionAuditor interface {
	CheckPermission(permissionID string, userRoles int64) bool
	CheckPermissionWithLevel(permissionID string, userRoles int64, userLevel int32) bool
	CheckRoute(area, controller, action string, userRoles int64) bool
	CheckRouteWithLevel(area, controller, action string, userRoles int64, userLevel int32) bool
	CheckRouteKeyWithLevel(routeKey string, userRoles int64, userLevel int32) bool
}

type permissionAuditor struct {
	routeProvider      IRouteProvider
	permissionProvider IPermissionProvider
	routes             map[string]*sproto.RouteDTO
	permissions        map[string]*sproto.PermissionDTO
}

func NewPermissionAuditor(permissionProvider IPermissionProvider, routeProvider IRouteProvider) IPermissionAuditor {
	r := new(permissionAuditor)
	r.permissionProvider = permissionProvider
	r.routeProvider = routeProvider
	err := r.ReloadRoutePermissions()
	u.LogFatal(err)
	return r
}

func (x *permissionAuditor) ReloadRoutePermissions() error {
	var err error

	if x.routeProvider != nil {
		x.routes, err = x.routeProvider.GetRoutes()
		if err != nil {
			return err
		}
	}

	if x.permissionProvider != nil {
		x.permissions, err = x.permissionProvider.GetPermissions()
		if err != nil {
			return err
		}
	}

	return nil
}

func (x *permissionAuditor) CheckPermission(permissionID string, userRoles int64) bool {
	return x.CheckPermissionWithLevel(permissionID, userRoles, 0)
}
func (x *permissionAuditor) CheckPermissionWithLevel(permissionID string, userRoles int64, userLevel int32) bool {
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
		return (permission.AllowedRoles&userRoles) > 0 && userLevel >= permission.Level
	}
}

func (x *permissionAuditor) CheckRoute(area, controller, action string, userRoles int64) bool {
	return x.CheckRouteWithLevel(area, controller, action, userRoles, 0)
}

func (x *permissionAuditor) CheckRouteWithLevel(area, controller, action string, userRoles int64, userLevel int32) bool {
	if x.routeProvider == nil {
		log.Warn("route provider is nil")
		return false
	}
	if x.permissionProvider == nil {
		log.Warn("permission provider is nil")
		return false
	}

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
		r := checkPermission(permission, userRoles, userLevel)
		if !r {
			log.Debugf("routeKey: %s_%s_%s, roles: %d, level: %d, permission: %v", area, controller, action, userRoles, userLevel, permission)
		}
		return r
	}

	log.Warnf("permission: %s does not exist", route.Permission_ID)
	return false
}

func (x *permissionAuditor) CheckRouteKeyWithLevel(routeKey string, userRoles int64, userLevel int32) bool {
	if x.routeProvider == nil {
		log.Warn("route provider is nil")
		return false
	}
	if x.permissionProvider == nil {
		log.Warn("permission provider is nil")
		return false
	}

	var route *sproto.RouteDTO
	var exists bool
	if route, exists = x.routes[routeKey]; !exists {
		log.Warnf("route: [%s] does not exist", routeKey)
		return false
	}

	if permission, exists := x.permissions[route.Permission_ID]; exists {
		r := checkPermission(permission, userRoles, userLevel)
		if !r {
			log.Debugf("routeKey: %s, roles: %d, level: %d, permission: %v", routeKey, userRoles, userLevel, permission)
		}
		return r
	}

	log.Warnf("permission: %s does not exist", route.Permission_ID)
	return false
}

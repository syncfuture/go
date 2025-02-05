package ssecurity

import (
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/sproto"
	"github.com/syncfuture/go/sslice"
	"github.com/syncfuture/go/u"
)

type IPermissionAuditor interface {
	CheckPermission(permissionID string, userRoles int64, userScopes []string) bool
	CheckPermissionWithLevel(permissionID string, userRoles int64, userLevel int32, userScopes []string) bool
	CheckRoute(area, controller, action string, userRoles int64, userScopes []string) bool
	CheckRouteWithLevel(area, controller, action string, userRoles int64, userLevel int32, userScopes []string) bool
	CheckRouteKeyWithLevel(routeKey string, userRoles int64, userLevel int32, userScopes []string) bool
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

func (x *permissionAuditor) CheckPermission(permissionID string, userRoles int64, userScopes []string) bool {
	return x.CheckPermissionWithLevel(permissionID, userRoles, 0, userScopes)
}
func (x *permissionAuditor) CheckPermissionWithLevel(permissionID string, userRoles int64, userLevel int32, userScopes []string) bool {
	if permission, exists := x.permissions[permissionID]; exists {
		return checkPermission(permission, userRoles, userLevel, userScopes)
	}

	log.Warnf("permission: %s does not exist", permissionID)
	return false
}

func checkPermission(permission *sproto.PermissionDTO, userRoles int64, userLevel int32, userScopes []string) bool {
	// If permission.AllowedScopes is limited, and (userScopes is empty or userScopes is not a subset of AllowedScopes), then return false
	if len(permission.AllowedScopes) > 0 && (len(userScopes) == 0 || !sslice.HasAllStr(permission.AllowedScopes, userScopes)) {
		return false
	}

	if permission.IsAllowGuest {
		return true
	} else if permission.IsAllowAnyUser {
		return userRoles > 0
	} else {
		return (permission.AllowedRoles&userRoles) > 0 && userLevel >= permission.Level
	}
}

func (x *permissionAuditor) CheckRoute(area, controller, action string, userRoles int64, userScopes []string) bool {
	return x.CheckRouteWithLevel(area, controller, action, userRoles, 0, userScopes)
}

func (x *permissionAuditor) CheckRouteWithLevel(area, controller, action string, userRoles int64, userLevel int32, userScopes []string) bool {
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
		r := checkPermission(permission, userRoles, userLevel, userScopes)
		if !r {
			log.Debugf("routeKey: %s_%s_%s, roles: %d, level: %d, permission: %v", area, controller, action, userRoles, userLevel, permission)
		}
		return r
	}

	log.Warnf("permission: %s does not exist", route.Permission_ID)
	return false
}

func (x *permissionAuditor) CheckRouteKeyWithLevel(routeKey string, userRoles int64, userLevel int32, userScopes []string) bool {
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
		r := checkPermission(permission, userRoles, userLevel, userScopes)
		if !r {
			log.Debugf("routeKey: %s, roles: %d, level: %d, permission: %v", routeKey, userRoles, userLevel, permission)
		}
		return r
	}

	log.Warnf("permission: %s does not exist", route.Permission_ID)
	return false
}

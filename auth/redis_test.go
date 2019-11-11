package auth

import (
	"testing"

	log "github.com/kataras/golog"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel("debug")
}

func TestGetRoutes(t *testing.T) {
	a := NewRedisRoutePermissionProvider("hubadmin", false, "Famous901", "localhost:6379")
	b, err := a.GetRoutes()
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
}

func TestGetPermissions(t *testing.T) {
	a := NewRedisRoutePermissionProvider("hubadmin", false, "Famous901", "localhost:6379")
	b, err := a.GetPermissions()
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
}

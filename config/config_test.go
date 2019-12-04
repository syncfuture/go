package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var _config IConfigProvider = NewJsonConfigProvider()

func TestGetString(t *testing.T) {
	r := _config.GetString("OIDC.ClientSecret")
	assert.Equal(t, "Famous901", r)
}

func TestGetBool(t *testing.T) {
	r := _config.GetBool("Dev.Debug")
	assert.Equal(t, true, r)
}

func TestGetInt(t *testing.T) {
	r := _config.GetInt("TestInt")
	assert.Equal(t, 79, r)
}

func TestGetStringSlice(t *testing.T) {
	r := _config.GetStringSlice("Redis.Addrs")
	assert.Len(t, r, 2)
}

func TestGetIntSlice(t *testing.T) {
	r := _config.GetIntSlice("TestIntSlice")
	assert.Len(t, r, 5)
}

func TestGetRedisConfig(t *testing.T) {
	r := _config.GetRedisConfig()
	assert.Equal(t, "Famous901", r.Password)
}

func TestGetOIDCConfig(t *testing.T) {
	r := _config.GetOIDCConfig()
	assert.Equal(t, "Famous901", r.ClientSecret)
}

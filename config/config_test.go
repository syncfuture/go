package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "xxxxxxxxxx", r.Password)
}

func TestGetOIDCConfig(t *testing.T) {
	r := _config.GetOIDCConfig()
	assert.Equal(t, "xxxxxxxxxx", r.ClientSecret)
}

func TestGetRabbitCConfig(t *testing.T) {
	r := _config.GetRabbitMQConfig()
	assert.Equal(t, "TestConsumer1", r.Nodes[0].Consumers[0].Name)
}

// func TestGetMapSlice(t *testing.T) {
// 	r := _config.GetMapSlice("Users")
// 	assert.Len(t, r, 1)
// }

// func TestGetMap(t *testing.T) {
// 	r := _config.GetMap("Redis")
// 	p := r.GetString("Password")
// 	assert.Equal(t, "xxxxxxxxxx", p)
// }

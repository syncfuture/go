package config

import (
	"github.com/syncfuture/go/soidc"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/go/rabbitmq"
)

type IConfigProvider interface {
	// GetMap(key string) MapConfiguration
	// GetMapSlice(key string) []MapConfiguration
	GetStruct(key string, target interface{}) error
	GetString(key string) string
	GetStringDefault(key string, defaultValue string) string
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetIntDefault(key string, defaultValue int) int
	GetStringSlice(key string) []string
	GetIntSlice(key string) []int

	GetRedisConfig() *sredis.RedisConfig
	GetOIDCConfig() *soidc.OIDCConfig
	GetRabbitMQConfig() *rabbitmq.RabbitMQConfig
}

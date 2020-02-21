package config

import (
	"reflect"

	"github.com/syncfuture/go/soidc"
	"github.com/syncfuture/go/sredis"

	"github.com/syncfuture/go/sjson"
)

type IConfigReader interface {
	ReadConfig(interface{}, ...string) error
}

type IConfigProvider interface {
	GetMap(key string) MapConfiguration
	GetMapSlice(key string) []MapConfiguration
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
}

type JsonConfigProvider struct {
	MapConfiguration
}

func NewJsonConfigProvider(args ...string) IConfigProvider {
	r := new(JsonConfigProvider)
	r.MapConfiguration = make(MapConfiguration)

	var configFile string
	if len(args) == 0 {
		configFile = "configs.json"
	} else {
		configFile = args[0]
	}

	r.ReadConfig(configFile, &r.MapConfiguration)

	return r
}

func (x *JsonConfigProvider) ReadConfig(configFile string, config interface{}) error {
	t := reflect.TypeOf(config)
	if t.Kind() != reflect.Ptr {
		panic("config must be a pointer")
	}

	err := sjson.DeserializeFromFile(configFile, config)
	return err
}

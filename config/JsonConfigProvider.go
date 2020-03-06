package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/syncfuture/go/u"

	"github.com/syncfuture/go/sjson"
	"github.com/syncfuture/go/soidc"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/go/rabbitmq"
)

type JsonConfigProvider struct {
	RawJson []byte
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

	configData, err := ioutil.ReadFile(configFile)
	u.LogFaltal(err)
	r.RawJson = configData
	err = json.Unmarshal(configData, &r.MapConfiguration)
	u.LogFaltal(err)

	return r
}

func (x *JsonConfigProvider) GetStruct(key string, target interface{}) error {
	return sjson.UnmarshalSection(x.RawJson, key, target)
}

func (x *JsonConfigProvider) GetRedisConfig() (r *sredis.RedisConfig) {
	err := x.GetStruct("Redis", &r)
	u.LogError(err)
	return r
}

func (x *JsonConfigProvider) GetOIDCConfig() (r *soidc.OIDCConfig) {
	err := x.GetStruct("OIDC", &r)
	u.LogError(err)
	return r
}

func (x *JsonConfigProvider) GetRabbitMQConfig() (r *rabbitmq.RabbitMQConfig) {
	err := x.GetStruct("RabbitMQ", &r)
	u.LogError(err)
	return r
}

package sconfig

import (
	"encoding/json"
	"os"

	"github.com/syncfuture/go/sjson"
	"github.com/syncfuture/go/u"
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

	configData, err := os.ReadFile(configFile)
	u.LogFaltal(err)
	r.RawJson = configData
	err = json.Unmarshal(configData, &r.MapConfiguration)
	u.LogFaltal(err)

	return r
}

func (x *JsonConfigProvider) GetStruct(key string, target interface{}) error {
	return sjson.UnmarshalSection(x.RawJson, key, target)
}

package config

import (
	"reflect"

	"github.com/syncfuture/go/json"
)

type IConfigProvider interface {
	GetConfig(interface{}, ...string) error
}

type JsonConfigProvider struct {
}

func CreateJsonConfigProvider() *JsonConfigProvider {
	return &JsonConfigProvider{}
}

func (x *JsonConfigProvider) GetConfig(configPtr interface{}, args ...string) error {
	if len(args) == 0 {
		panic("must pass a config file path")
	}

	t := reflect.TypeOf(configPtr)
	if t.Kind() != reflect.Ptr {
		panic("configPtr must be a pointer")
	}

	err := json.DeserializeFromFile(args[0], configPtr)
	return err
}

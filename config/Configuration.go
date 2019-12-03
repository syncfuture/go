package config

import (
	log "github.com/kataras/golog"
	"strings"
)

type Configuration map[string]interface{}

func (x *Configuration) GetString(key string) string {
	v := getValue(key, *x)
	if v != nil {
		r, ok := v.(string)
		if !ok {
			log.Warnf("convert failed. %t -> string", v)
		} else {
			return r
		}
	}
	return ""
}

func (x *Configuration) GetBool(key string) bool {
	v := getValue(key, *x)
	if v != nil {
		r, ok := v.(bool)
		if !ok {
			log.Warnf("convert failed. %t -> bool", v)
		} else {
			return r
		}
	}
	return false
}

func (x *Configuration) GetFloat64(key string) float64 {
	v := getValue(key, *x)
	if v != nil {
		r, ok := v.(float64)
		if !ok {
			log.Warnf("convert failed. %t -> float64", v)
		} else {
			return r
		}
	}
	return 0
}

func (x *Configuration) GetInt(key string) int {
	r := x.GetFloat64(key)
	return int(r)
}

func getValue(key string, c Configuration) interface{} {
	keys := strings.Split(key, ".")
	keyCount := len(keys)

	for i := 0; i < keyCount; i++ {
		if i < keyCount-1 {
			c = getNode(keys[i], c)
			if c == nil {
				break
			}
		} else {
			r, ok := c[keys[i]]
			if ok {
				return r
			}
		}
	}

	return nil
}

func getNode(key string, c map[string]interface{}) map[string]interface{} {
	v, ok := c[key]
	if ok {
		return v.(map[string]interface{})
	}
	return nil
}

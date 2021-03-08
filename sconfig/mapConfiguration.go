package sconfig

import (
	"strings"

	log "github.com/kataras/golog"
)

type MapConfiguration map[string]interface{}

func (x *MapConfiguration) GetMap(key string) (r MapConfiguration) {
	v := getValue(key, *x)
	if v != nil {
		a, ok := v.(interface{})
		if !ok {
			log.Warnf("convert failed. %t -> interface{}", v)
		} else {
			b := a.(map[string]interface{})
			r = MapConfiguration(b)

			return r
		}
	}
	return nil
}

func (x *MapConfiguration) GetMapSlice(key string) (r []MapConfiguration) {
	r = make([]MapConfiguration, 0)
	v := getValue(key, *x)
	if v != nil {
		a, ok := v.([]interface{})
		if !ok {
			log.Warnf("convert failed. %t -> []interface{}", v)
		} else {
			for _, b := range a {
				c := b.(map[string]interface{})
				r = append(r, MapConfiguration(c))
			}

			return r
		}
	}
	return nil
}

func (x *MapConfiguration) GetString(key string) string {
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

func (x *MapConfiguration) GetStringDefault(key, defaultValue string) string {
	r := x.GetString(key)
	if r != "" {
		return r
	}
	return defaultValue
}

func (x *MapConfiguration) GetBool(key string) bool {
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

func (x *MapConfiguration) GetFloat64(key string) float64 {
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

func (x *MapConfiguration) GetInt(key string) int {
	r := x.GetFloat64(key)
	return int(r)
}
func (x *MapConfiguration) GetIntDefault(key string, defaultValue int) int {
	r := x.GetInt(key)
	if r != 0 {
		return r
	}
	return defaultValue
}

func (x *MapConfiguration) GetStringSlice(key string) []string {
	v := getValue(key, *x)
	if v != nil {
		slice, ok := v.([]interface{})
		if !ok {
			log.Warnf("convert failed. %t -> interface slice", v)
		} else {
			var r []string
			for _, e := range slice {
				a, ok := e.(string)
				if ok {
					r = append(r, a)
				}
			}
			return r
		}
	}
	return make([]string, 0)
}

func (x *MapConfiguration) GetIntSlice(key string) []int {
	v := getValue(key, *x)
	if v != nil {
		slice, ok := v.([]interface{})
		if !ok {
			log.Warnf("convert failed. %t -> interface slice", v)
		} else {
			var r []int
			for _, e := range slice {
				a, ok := e.(float64)
				if ok {
					r = append(r, int(a))
				}
			}
			return r
		}
	}
	return make([]int, 0)
}

func getValue(key string, c MapConfiguration) interface{} {
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

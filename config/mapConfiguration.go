package config

import (
	log "github.com/kataras/golog"
	"github.com/syncfuture/go/soidc"
	"github.com/syncfuture/go/sredis"
	"strings"
)

type mapConfiguration map[string]interface{}

func (x *mapConfiguration) GetString(key string) string {
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

func (x *mapConfiguration) GetStringDefault(key, defaultValue string) string {
	r := x.GetString(key)
	if r != "" {
		return r
	}
	return defaultValue
}

func (x *mapConfiguration) GetBool(key string) bool {
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

func (x *mapConfiguration) GetFloat64(key string) float64 {
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

func (x *mapConfiguration) GetInt(key string) int {
	r := x.GetFloat64(key)
	return int(r)
}
func (x *mapConfiguration) GetIntDefault(key string, defaultValue int) int {
	r := x.GetInt(key)
	if r != 0 {
		return r
	}
	return defaultValue
}

func (x *mapConfiguration) GetStringSlice(key string) []string {
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

func (x *mapConfiguration) GetIntSlice(key string) []int {
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

func (x *mapConfiguration) GetRedisConfig() *sredis.RedisConfig {
	r := new(sredis.RedisConfig)
	r.Addrs = x.GetStringSlice("Redis.Addrs")
	r.Password = x.GetString("Redis.Password")
	r.ClusterEnabled = x.GetBool("Redis.ClusterEnabled")
	return r
}

func (x *mapConfiguration) GetOIDCConfig() *soidc.OIDCConfig {
	r := new(soidc.OIDCConfig)
	r.ClientID = x.GetString("OIDC.Addrs")
	r.ClientSecret = x.GetString("OIDC.ClientSecret")
	r.PassportURL = x.GetString("OIDC.PassportURL")
	r.JWKSURL = x.GetString("OIDC.JWKSURL")
	r.SignInCallbackURL = x.GetString("OIDC.SignInCallbackURL")
	r.SignOutCallbackURL = x.GetString("OIDC.SignOutCallbackURL")
	r.AccessDeniedURL = x.GetString("OIDC.AccessDeniedURL")
	r.Scopes = x.GetStringSlice("OIDC.Scopes")
	return r
}

func getValue(key string, c mapConfiguration) interface{} {
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

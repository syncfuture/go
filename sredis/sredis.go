package sredis

import (
	"github.com/go-redis/redis"
	log "github.com/kataras/golog"
)

type RedisConfig struct {
	Addrs          []string
	Password       string
	ClusterEnabled bool
}

func NewClient(config *RedisConfig) redis.Cmdable {
	if len(config.Addrs) == 0 {
		log.Fatal("addrs cannot be empty")
	}

	if config.ClusterEnabled {
		c := &redis.ClusterOptions{
			Addrs: config.Addrs,
		}
		if config.Password != "" {
			c.Password = config.Password
		}
		return redis.NewClusterClient(c)
	} else {
		c := &redis.Options{
			Addr: config.Addrs[0],
		}
		if config.Password != "" {
			c.Password = config.Password
		}
		return redis.NewClient(c)
	}
}

package goredis

import (
	"github.com/go-redis/redis"
	log "github.com/kataras/golog"
)

func NewClient(clusterEnabled bool, password string, addrs ...string) redis.Cmdable {
	if len(addrs) == 0 {
		log.Fatal("addrs cannot be empty")
	}

	if clusterEnabled {
		c := &redis.ClusterOptions{
			Addrs: addrs,
		}
		if password != "" {
			c.Password = password
		}
		return redis.NewClusterClient(c)
	} else {
		c := &redis.Options{
			Addr: addrs[0],
		}
		if password != "" {
			c.Password = password
		}
		return redis.NewClient(c)
	}
}

package sredis

import (
	"github.com/go-redis/redis/v7"
	log "github.com/kataras/golog"
	u "github.com/syncfuture/go/util"
)

type RedisConfig struct {
	Addrs    []string
	Password string
}

func NewClient(config *RedisConfig) redis.Cmdable {
	addrCount := len(config.Addrs)
	if addrCount == 0 {
		log.Fatal("addrs cannot be empty")
		return nil
	} else if addrCount == 1 {
		c := &redis.Options{
			Addr: config.Addrs[0],
		}
		if config.Password != "" {
			c.Password = config.Password
		}
		return redis.NewClient(c)
	} else {
		c := &redis.ClusterOptions{
			Addrs: config.Addrs,
		}
		if config.Password != "" {
			c.Password = config.Password
		}
		return redis.NewClusterClient(c)
	}
}

func GetPagedKeys(client redis.Cmdable, match string, pageSize int64) (cursor uint64, r []string, err error) {
	r, cursor, err = client.Scan(cursor, match, pageSize).Result()
	u.LogError(err)
	return
}

func GetAllKeys(client redis.Cmdable, match string, pageSize int64) (r []string, err error) {
	var cursor uint64
	for {
		var ks []string
		cursor, ks, err = GetPagedKeys(client, match, pageSize)
		r = append(r, ks...)
		if cursor == 0 {
			break
		}
	}

	return
}

package sredis

import (
	"context"
	"sort"

	"github.com/redis/go-redis/v9"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
)

type RedisConfig struct {
	Addrs    []string
	Password string
	DB       int
	// ClusterEnabled bool
}

func NewClient(config *RedisConfig) redis.UniversalClient {
	addrCount := len(config.Addrs)
	if addrCount == 0 {
		log.Fatal("addrs cannot be empty")
		return nil
		// } else if addrCount == 1 && !config.ClusterEnabled {
	} else if addrCount == 1 {
		c := &redis.Options{
			Addr: config.Addrs[0],
			DB:   config.DB,
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

func GetPagedKeys(client redis.Cmdable, cursor uint64, match string, pageSize int64) (uint64, []string) {
	r, newCursor, err := client.Scan(context.Background(), cursor, match, pageSize).Result()
	u.LogError(err)
	return newCursor, r
}

func GetAllKeys(client redis.Cmdable, match string, pageSize int64) (r []string) {
	var cursor uint64
	for {
		var ks []string
		cursor, ks = GetPagedKeys(client, cursor, match, pageSize)
		r = append(r, ks...)
		if cursor == 0 {
			break
		}
	}

	sort.Strings(r)

	return
}

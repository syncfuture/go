package surl

import (
	"github.com/go-redis/redis"
	"github.com/syncfuture/go/sredis"
	u "github.com/syncfuture/go/util"
)

const (
	key = "ecp:URIS"
)

type redisURLProvider struct {
	client redis.Cmdable
}

func NewRedisURLProvider(redisConfig *sredis.RedisConfig) IURLProvider {
	return &redisURLProvider{
		client: sredis.NewClient(redisConfig),
	}
}

func (x *redisURLProvider) GetURL(urlKey string) string {
	cmd := x.client.HGet(key, urlKey)
	r, err := cmd.Result()
	u.LogError(err)
	return r
}

package surl

import (
	"regexp"
	"time"

	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
	"github.com/syncfuture/go/sredis"
	u "github.com/syncfuture/go/util"
)

const (
	key = "ecp:URIS"
)

var (
	_cache = cache.New(1*time.Hour, 1*time.Hour)
	_regex = regexp.MustCompile("{{[^}]+}}")
)

type redisURLProvider struct {
	client redis.Cmdable
}

// NewRedisURLProvider create new url provider
func NewRedisURLProvider(redisConfig *sredis.RedisConfig) IURLProvider {
	return &redisURLProvider{
		client: sredis.NewClient(redisConfig),
	}
}

// GetURL get url from redis
func (x *redisURLProvider) GetURL(urlKey string) string {
	cmd := x.client.HGet(key, urlKey)
	r, err := cmd.Result()
	u.LogError(err)
	return r
}

// GetURLCache get url from cache, if does not exist, call GetURL and put it into cache
func (x *redisURLProvider) GetURLCache(urlKey string) string {
	r, found := _cache.Get(urlKey)
	if !found {
		r = x.GetURL(urlKey)
		_cache.Set(urlKey, r, 1*time.Hour)
	}
	return r.(string)
}

func (x *redisURLProvider) RenderURL(url string) string {
	return _regex.ReplaceAllStringFunc(url, x.GetURL)
}

func (x *redisURLProvider) RenderURLCache(url string) string {
	return _regex.ReplaceAllStringFunc(url, func(match string) string {
		return x.GetURLCache(match[2 : len(match)-2])
	})
}

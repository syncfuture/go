package surl

import (
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/patrickmn/go-cache"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/go/u"
)

const (
	key = "ecp:URIS"
)

var (
	_cache = cache.New(1*time.Hour, 1*time.Hour)
	_regex = regexp.MustCompile("{{URI '([^']+)'}}")
)

type redisURLProvider struct {
	client redis.Cmdable
}

// NewRedisURLProvider create new url provider
func NewRedisURLProvider(redisConfig *sredis.RedisConfig) IURLProvider {
	if len(redisConfig.Addrs) == 0 || redisConfig.Addrs[0] == "" {
		log.Fatal("cannot find 'Redis.Addrs' config")
	}
	return &redisURLProvider{
		client: sredis.NewClient(redisConfig),
	}
}

// GetURL get url from redis
func (x *redisURLProvider) GetURL(urlKey string) string {
	if urlKey == "" {
		return urlKey
	}

	cmd := x.client.HGet(key, urlKey)
	r, err := cmd.Result()
	u.LogError(err)
	return r
}

// GetURLCache get url from cache, if does not exist, call GetURL and put it into cache
func (x *redisURLProvider) GetURLCache(urlKey string) string {
	if urlKey == "" {
		return urlKey
	}

	r, found := _cache.Get(urlKey)
	if !found {
		r = x.GetURL(urlKey)
		_cache.Set(urlKey, r, 1*time.Hour)
	}
	return r.(string)
}

func (x *redisURLProvider) RenderURL(url string) string {
	if url == "" {
		return url
	}

	// return _regex.ReplaceAllStringFunc(url, x.GetURL)
	matches := _regex.FindAllStringSubmatch(url, 5)
	if len(matches) > 0 && len(matches[0]) > 1 {
		o := matches[0][0]
		n := x.GetURL(matches[0][1])
		return strings.ReplaceAll(url, o, n)
	}
	return url
}

func (x *redisURLProvider) RenderURLCache(url string) string {
	if url == "" {
		return url
	}

	r, found := _cache.Get(url)
	if !found {
		r = x.RenderURL(url)
		_cache.Set(url, r, 1*time.Hour)
	}
	return r.(string)
}

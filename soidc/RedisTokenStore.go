package soidc

import (
	"github.com/go-redis/redis"
	"github.com/syncfuture/go/json"
	"github.com/syncfuture/go/sredis"
	u "github.com/syncfuture/go/util"
	"golang.org/x/oauth2"
)

type RedisTokenStore struct {
	redis redis.Cmdable
	Key   string
}

func NewRedisTokenStore(projectName string, config *sredis.RedisConfig) ITokenStore {
	r := new(RedisTokenStore)

	r.redis = sredis.NewClient(config)

	r.Key = key

	return r
}

func (x *RedisTokenStore) Save(userID string, token *oauth2.Token) error {
	j, err := json.Serialize(token)
	if u.LogError(err) {
		return err
	}

	// Todo: 加密
	cmd := x.redis.HSet(x.Key, userID, j)
	return cmd.Err()
}

func (x *RedisTokenStore) Get(userID string) (*oauth2.Token, error) {
	cmd := x.redis.HGet(x.Key, userID)
	j, err := cmd.Result()
	if u.LogError(err) {
		return nil, err
	}

	// Todo: 解密

	r := new(oauth2.Token)
	err = json.Deserialize(j, r)
	if u.LogError(err) {
		return nil, err
	}

	return r, nil
}

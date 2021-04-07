package ssecurity

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/syncfuture/go/serr"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/sproto"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/go/u"
)

type RedisRouteProvider struct {
	redis    redis.Cmdable
	RouteKey string
}

func NewRedisRouteProvider(routeKey string, config *sredis.RedisConfig) IRouteProvider {
	if routeKey == "" {
		log.Fatal("routeKey cannot be empty")
	}

	r := new(RedisRouteProvider)

	r.redis = sredis.NewClient(config)

	r.RouteKey = routeKey

	return r
}

// *******************************************************************************************************************************
// Route
func (x *RedisRouteProvider) CreateRoute(in *sproto.RouteDTO) error {
	j, err := json.Marshal(in)
	if err != nil {
		return serr.WithStack(err)
	}

	cmd := x.redis.HSet(context.Background(), x.RouteKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisRouteProvider) GetRoute(id string) (*sproto.RouteDTO, error) {
	cmd := x.redis.HGet(context.Background(), x.RouteKey, id)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}

	r := new(sproto.RouteDTO)
	j, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(j), r)
	u.LogError(err)
	return r, err
}
func (x *RedisRouteProvider) UpdateRoute(in *sproto.RouteDTO) error {
	j, err := json.Marshal(in)
	if err != nil {
		return serr.WithStack(err)
	}

	cmd := x.redis.HSet(context.Background(), x.RouteKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisRouteProvider) RemoveRoute(id string) error {
	cmd := x.redis.HDel(context.Background(), x.RouteKey, id)
	return cmd.Err()
}
func (x *RedisRouteProvider) GetRoutes() (map[string]*sproto.RouteDTO, error) {
	cmd := x.redis.HGetAll(context.Background(), x.RouteKey)
	r, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	m := make(map[string]*sproto.RouteDTO, len(r))
	for key, value := range r {
		dto := new(sproto.RouteDTO)
		err = json.Unmarshal([]byte(value), dto)
		if !u.LogError(err) {
			m[key] = dto
		}
	}
	return m, err
}

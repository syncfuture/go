package surl

import (
	"testing"

	"github.com/syncfuture/go/sredis"
)

func TestRenderURL(t *testing.T) {
	a := NewRedisURLProvider(&sredis.RedisConfig{
		Addrs:    []string{"localhost:6379"},
		Password: "Famous901",
	})

	b := a.RenderURLCache("{{hub}}/v1/product/items")
	t.Log(b)
}

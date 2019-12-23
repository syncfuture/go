package surl

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/syncfuture/go/sredis"
)

func TestRenderURL(t *testing.T) {
	a := NewRedisURLProvider(&sredis.RedisConfig{
		Addrs:    []string{"localhost:6379"},
		Password: "Famous901",
	})

	b := a.RenderURLCache("{{URI 'hubapi'}}/product/items")
	assert.Equal(t, "https://di.dreamvat.com/product/items", b)
	t.Log(b)
}

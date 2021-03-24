package spool

import (
	"net/http"
	"sync"
	"time"
)

type ICookiePool interface {
	GetCookie() *http.Cookie
	PutCookie(*http.Cookie)
}

type syncCookiePool struct {
	pool *sync.Pool
}

func NewSyncCookiePool() ICookiePool {
	var newPool syncCookiePool

	newPool.pool = &sync.Pool{
		New: func() interface{} {
			return new(http.Cookie)
		},
	}

	return &newPool
}

func (x *syncCookiePool) GetCookie() *http.Cookie {
	return x.pool.Get().(*http.Cookie)
}

func (x *syncCookiePool) PutCookie(c *http.Cookie) {
	if c != nil {
		c.Domain = ""
		c.Expires = time.Time{}
		c.HttpOnly = false
		c.MaxAge = 0
		c.Name = ""
		c.Path = ""
		c.Raw = ""
		c.RawExpires = ""
		c.SameSite = 0
		c.Secure = false
		c.Value = ""
		c.Unparsed = nil
	}
	x.pool.Put(c)
}

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
	r := x.pool.Get().(*http.Cookie)
	r.Domain = ""
	r.Expires = time.Time{}
	r.HttpOnly = false
	r.MaxAge = 0
	r.Name = ""
	r.Path = ""
	r.Raw = ""
	r.RawExpires = ""
	r.SameSite = 0
	r.Secure = false
	r.Value = ""
	r.Unparsed = nil
	return r
}

func (x *syncCookiePool) PutCookie(c *http.Cookie) {
	x.pool.Put(c)
}

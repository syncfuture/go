package spool

import (
	"sync"
)

type IMapPool interface {
	GetMap() *map[string]interface{}
	PutMap(*map[string]interface{})
}

type syncMapPool struct {
	pool *sync.Pool
}

func NewSyncMapPool() IMapPool {
	var newPool syncMapPool

	newPool.pool = &sync.Pool{
		New: func() interface{} {
			r := make(map[string]interface{})
			return &r
		},
	}

	return &newPool
}

func (x *syncMapPool) GetMap() *map[string]interface{} {
	r := x.pool.Get().(*map[string]interface{})
	for k := range *r {
		delete(*r, k)
	}
	return r
}

func (x *syncMapPool) PutMap(m *map[string]interface{}) {
	x.pool.Put(m)
}

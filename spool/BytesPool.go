package spool

import (
	"sync"
)

type IBytesPool interface {
	GetBytes() *[]byte
	PutBytes(*[]byte)
}

type syncBytesPool struct {
	pool *sync.Pool
}

func NewSyncBytesPool(buf_size int) IBytesPool {
	return &syncBytesPool{
		pool: &sync.Pool{
			New: func() interface{} {
				b := make([]byte, buf_size)
				return &b
			},
		},
	}
}

func (bp *syncBytesPool) GetBytes() *[]byte {
	r := bp.pool.Get().(*[]byte)
	for i := range *r {
		(*r)[i] = 0
	}
	return r
}

func (bp *syncBytesPool) PutBytes(b *[]byte) {
	bp.pool.Put(b)
}

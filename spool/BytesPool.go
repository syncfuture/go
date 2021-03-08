package spool

import (
	"sync"
)

type BytesPool interface {
	GetBytes() *[]byte
	PutBytes(*[]byte)
}

type syncBytesPool struct {
	pool      *sync.Pool
	makeBytes func() interface{}
}

func NewSyncBytesPool(buf_size int) BytesPool {
	var newPool syncBytesPool

	newPool.makeBytes = func() interface{} {
		b := make([]byte, 1024)
		return &b
	}
	newPool.pool = &sync.Pool{}
	newPool.pool.New = newPool.makeBytes

	return &newPool
}

func (bp *syncBytesPool) GetBytes() (b *[]byte) {
	pool_object := bp.pool.Get()

	b, ok := pool_object.(*[]byte)
	if !ok { // explicitly make buffer if sync.Pool returns nil:
		b = bp.makeBytes().(*[]byte)
	}
	return
}

func (bp *syncBytesPool) PutBytes(b *[]byte) {
	bp.pool.Put(b)
}

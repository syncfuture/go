package spool

import (
	"bytes"
	"sync"
)

type IBufferPool interface {
	GetBuffer() *bytes.Buffer
	PutBuffer(*bytes.Buffer)
}

type syncBufferPool struct {
	pool *sync.Pool
}

func NewSyncBufferPool(buf_size int) IBufferPool {
	var newPool syncBufferPool

	newPool.pool = &sync.Pool{
		New: func() interface{} {
			var b bytes.Buffer
			b.Grow(buf_size)
			return &b
		},
	}

	return &newPool
}

func (x *syncBufferPool) GetBuffer() *bytes.Buffer {
	b := x.pool.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

func (bp *syncBufferPool) PutBuffer(b *bytes.Buffer) {
	bp.pool.Put(b)
}

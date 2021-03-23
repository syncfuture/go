package spool

import (
	"bytes"
	"sync"
)

type BufferPool interface {
	GetBuffer() *bytes.Buffer
	PutBuffer(*bytes.Buffer)
}

type syncBufferPool struct {
	pool *sync.Pool
	// makeBuffer func() interface{}
}

func NewSyncBufferPool(buf_size int) BufferPool {
	var newPool syncBufferPool

	// newPool.makeBuffer = func() interface{} {
	// 	var b bytes.Buffer
	// 	b.Grow(buf_size)
	// 	return &b
	// }
	newPool.pool = &sync.Pool{
		New: func() interface{} {
			var b bytes.Buffer
			b.Grow(buf_size)
			return &b
		},
	}
	// newPool.pool.New = newPool.makeBuffer

	return &newPool
}

func (x *syncBufferPool) GetBuffer() (b *bytes.Buffer) {
	// pool_object := bp.pool.Get()
	// b, ok := pool_object.(*bytes.Buffer)
	// if !ok { // explicitly make buffer if sync.Pool returns nil:
	// 	b = bp.makeBuffer().(*bytes.Buffer)
	// }
	// } else {
	// 	b.Reset()
	// }

	return x.pool.Get().(*bytes.Buffer)
}

func (bp *syncBufferPool) PutBuffer(b *bytes.Buffer) {
	if b != nil {
		b.Reset()
	}
	bp.pool.Put(b)
}

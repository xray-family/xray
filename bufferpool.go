package uRouter

import (
	"bytes"
	"sync"
)

const defaultBufferSize = 1024

type BufferPool interface {
	SetSize(size int)
	Get() *bytes.Buffer
	Put(b *bytes.Buffer)
}

func DefaultBufferPool() BufferPool {
	return defaultBufferPool
}

func newBufferPool() *bufferPool {
	bp := &bufferPool{
		size: defaultBufferSize,
	}
	bp.p.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, bp.size))
	}
	return bp
}

type bufferPool struct {
	p    sync.Pool
	size int
}

func (c *bufferPool) SetSize(size int) {
	c.size = size
}

func (c *bufferPool) Get() *bytes.Buffer {
	b := c.p.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

func (c *bufferPool) Put(b *bytes.Buffer) {
	c.p.Put(b)
}

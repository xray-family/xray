package uRouter

import (
	"bytes"
	"sync"
)

const defaultBufferSize = 4 * 1024

type BufferPool interface {
	Get() *bytes.Buffer
	Put(b *bytes.Buffer)
}

func DefaultBufferPool() BufferPool {
	return defaultBufferPool
}

func newBufferPool() *bufferPool {
	return &bufferPool{
		p: sync.Pool{New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, defaultBufferSize))
		}},
	}
}

type bufferPool struct {
	p sync.Pool
}

func (c *bufferPool) Get() *bytes.Buffer {
	b := c.p.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

func (c *bufferPool) Put(b *bytes.Buffer) {
	c.p.Put(b)
}

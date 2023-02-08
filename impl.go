package uRouter

import (
	"bytes"
	"log"
	"sync"
)

const defaultBufferSize = 1024

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

type logger struct{}

func (c *logger) Debug(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (c *logger) Warn(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (c *logger) Error(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (c *logger) Info(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (c *logger) Panic(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

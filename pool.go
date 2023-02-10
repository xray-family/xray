package uRouter

import (
	"bytes"
	"github.com/lxzan/uRouter/constant"
	"sync"
)

func newBufferPool() *bufferPool {
	bp := &bufferPool{
		p1:  &sync.Pool{},
		p4:  &sync.Pool{},
		p8:  &sync.Pool{},
		p16: &sync.Pool{},
	}
	bp.p1.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, constant.BufferLeveL1))
	}
	bp.p4.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, constant.BufferLeveL4))
	}
	bp.p8.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, constant.BufferLeveL8))
	}
	bp.p16.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, constant.BufferLeveL16))
	}
	return bp
}

type bufferPool struct {
	p1  *sync.Pool
	p4  *sync.Pool
	p8  *sync.Pool
	p16 *sync.Pool
}

func (c *bufferPool) Get(n int) *bytes.Buffer {
	var b *bytes.Buffer
	if n <= constant.BufferLeveL1 {
		b = c.p1.Get().(*bytes.Buffer)
	} else if n <= constant.BufferLeveL4 {
		b = c.p4.Get().(*bytes.Buffer)
	} else if n <= constant.BufferLeveL8 {
		b = c.p8.Get().(*bytes.Buffer)
	} else if n <= constant.BufferLeveL16 {
		b = c.p16.Get().(*bytes.Buffer)
	} else {
		b = bytes.NewBuffer(make([]byte, 0, n))
	}
	return b
}

func (c *bufferPool) Put(b *bytes.Buffer) {
	if b == nil {
		return
	}

	n := b.Len()
	b.Reset()
	if n <= constant.BufferLeveL1 {
		c.p1.Put(b)
	} else if n <= constant.BufferLeveL4 {
		c.p4.Put(b)
	} else if n <= constant.BufferLeveL8 {
		c.p8.Put(b)
	} else if n <= constant.BufferLeveL16 {
		c.p16.Put(b)
	}
}

package xray

import (
	"bytes"
	"github.com/lxzan/xray/constant"
	"sync"
)

func newBufferPool() *bufferPool {
	bp := &bufferPool{
		p1: &sync.Pool{},
		p2: &sync.Pool{},
		p4: &sync.Pool{},
	}
	bp.p1.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, constant.BufferLeveL1))
	}
	bp.p2.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, constant.BufferLeveL2))
	}
	bp.p4.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, constant.BufferLeveL4))
	}
	return bp
}

type bufferPool struct {
	p1 *sync.Pool
	p2 *sync.Pool
	p4 *sync.Pool
}

func (c *bufferPool) Get(n int) *bytes.Buffer {
	var b *bytes.Buffer
	if n <= constant.BufferLeveL1 {
		b = c.p1.Get().(*bytes.Buffer)
	} else if n <= constant.BufferLeveL2 {
		b = c.p2.Get().(*bytes.Buffer)
	} else if n <= constant.BufferLeveL4 {
		b = c.p4.Get().(*bytes.Buffer)
	} else {
		b = bytes.NewBuffer(make([]byte, 0, n))
	}
	return b
}

func (c *bufferPool) Put(b *bytes.Buffer) {
	if b == nil {
		return
	}

	n := b.Cap()
	b.Reset()
	if n <= constant.BufferLeveL1 {
		c.p1.Put(b)
	} else if n <= constant.BufferLeveL2 {
		c.p2.Put(b)
	} else if n <= constant.BufferLeveL4 {
		c.p4.Put(b)
	}
}

func HeaderPool() *headerPool {
	return defaultHeaderPool
}

func newHeaderPool() *headerPool {
	var c = new(headerPool)
	c.Register(constant.IdMapHeader, func() Header {
		return &MapHeader{}
	})
	return c
}

type headerPool struct {
	pools [8]*sync.Pool
}

func (c *headerPool) Register(id int, generate func() Header) {
	c.pools[id] = &sync.Pool{New: func() interface{} {
		return generate()
	}}
}

func (c *headerPool) Get(id uint8) Header {
	return c.pools[id].Get().(Header)
}

func (c *headerPool) Put(id uint8, h Header) {
	c.pools[id].Put(h)
}

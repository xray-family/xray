package uRouter

import (
	"github.com/lxzan/uRouter/internal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBufferPool(t *testing.T) {
	var as = assert.New(t)
	p := BufferPool()
	p = newBufferPool()

	b0 := p.Get(512)
	as.Equal(internal.BufferLeveL1, b0.Cap())
	as.Equal(0, b0.Len())

	b1 := p.Get(1200)
	as.Equal(internal.BufferLeveL2, b1.Cap())
	as.Equal(0, b1.Len())

	b2 := p.Get(3200)
	as.Equal(internal.BufferLeveL4, b2.Cap())
	as.Equal(0, b2.Len())

	b5 := p.Get(19200)
	as.Equal(19200, b5.Cap())
	as.Equal(0, b5.Len())

	p.Put(b0)
	p.Put(b1)
	p.Put(b2)
	p.Put(b5)
	p.Put(nil)

	b0 = p.Get(512)
	as.Equal(internal.BufferLeveL1, b0.Cap())
	as.Equal(0, b0.Len())

	b1 = p.Get(1200)
	as.Equal(internal.BufferLeveL2, b1.Cap())
	as.Equal(0, b1.Len())

	b2 = p.Get(3200)
	as.Equal(internal.BufferLeveL4, b2.Cap())
	as.Equal(0, b2.Len())

	b5 = p.Get(19200)
	as.Equal(19200, b5.Cap())
	as.Equal(0, b5.Len())
}

func TestHeaderPool(t *testing.T) {
	as := assert.New(t)

	p := HeaderPool()
	p = newHeaderPool()
	h1 := p.Get(internal.MapHeaderNumber)
	_, ok1 := h1.(*MapHeader)
	as.Equal(true, ok1)
	as.Equal(0, h1.Len())

	defaultHeaderPool.Put(&HttpHeader{Header: http.Header{}})

	h2 := defaultHeaderPool.Get(internal.MapHeaderNumber)
	_, ok2 := h2.(*MapHeader)
	as.Equal(true, ok2)
	as.Equal(0, h2.Len())

	h3, ok3 := newHeaderPool().Get(internal.HttpHeaderNumber).(HttpHeader)
	as.Equal(true, ok3)
	as.Equal(0, h3.Len())
}

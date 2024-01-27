package xray

import (
	"github.com/lxzan/xray/constant"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBufferPool(t *testing.T) {
	var as = assert.New(t)
	p := BufferPool()
	p = newBufferPool()

	b0 := p.Get(512)
	as.Equal(constant.BufferLeveL1, b0.Cap())
	as.Equal(0, b0.Len())

	b1 := p.Get(1200)
	as.Equal(constant.BufferLeveL2, b1.Cap())
	as.Equal(0, b1.Len())

	b2 := p.Get(3200)
	as.Equal(constant.BufferLeveL4, b2.Cap())
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
	as.Equal(constant.BufferLeveL1, b0.Cap())
	as.Equal(0, b0.Len())

	b1 = p.Get(1200)
	as.Equal(constant.BufferLeveL2, b1.Cap())
	as.Equal(0, b1.Len())

	b2 = p.Get(3200)
	as.Equal(constant.BufferLeveL4, b2.Cap())
	as.Equal(0, b2.Len())

	b5 = p.Get(19200)
	as.Equal(19200, b5.Cap())
	as.Equal(0, b5.Len())
}

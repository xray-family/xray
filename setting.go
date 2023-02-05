package uRouter

import (
	"bytes"
	"sync"
)

var (
	defaultJsonCodec Codec = new(stdJsonCodec)

	defaultGenerator = func() Header {
		return &MapHeader{}
	}

	StdJsonCodec Codec = new(stdJsonCodec)

	// TextHeader 文本类型头部编码, 4字节, 最大长度=9999
	TextHeader *HeaderCodec

	// BinaryHeader 二进制类型头部编码, 2字节, 最大长度=65535
	BinaryHeader *HeaderCodec

	DefaultBufferPool = newBufferPool()
)

func init() {
	TextHeader = NewHeaderCodec(TextLengthEncoding, defaultJsonCodec, defaultGenerator)

	BinaryHeader = NewHeaderCodec(BinaryLengthEncoding, defaultJsonCodec, defaultGenerator)
}

// SetJsonCodec 设置JSON编码器
func SetJsonCodec(codec Codec) {
	defaultJsonCodec = codec

	TextHeader = NewHeaderCodec(TextLengthEncoding, defaultJsonCodec, func() Header {
		return &MapHeader{}
	})

	BinaryHeader = NewHeaderCodec(BinaryLengthEncoding, defaultJsonCodec, func() Header {
		return &MapHeader{}
	})
}

const BufferSize = 4 * 1024

func newBufferPool() *bufferPool {
	return &bufferPool{
		p: sync.Pool{New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, BufferSize))
		}},
	}
}

type bufferPool struct {
	p sync.Pool
}

func (c *bufferPool) Get() *bytes.Buffer {
	return c.p.Get().(*bytes.Buffer)
}

func (c *bufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	c.p.Put(b)
}

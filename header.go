package uRouter

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/lxzan/uRouter/internal"
	"math"
	"net/http"
	"strconv"
)

const (
	XPath       = "X-Path"
	XRealIP     = "X-Real-Ip"
	ContentType = "Content-Type"

	MimeJson   = "application/json; charset=utf-8"
	MimeStream = "application/octet-stream"

	binaryLengthEncoding headerLengthEncoding = 2
	textLengthEncoding   headerLengthEncoding = 4
)

var (
	errHeaderTooLarge = errors.New("header size too large")
	errHeaderSize     = errors.New("header size error")
)

type (
	HeaderCodec struct {
		codec          Codec
		lengthEncoding headerLengthEncoding
		generator      func() Header
	}

	headerLengthEncoding uint8
)

func (c headerLengthEncoding) MaxLength() int {
	if c == binaryLengthEncoding {
		return math.MaxUint16
	}
	return 1e4 - 1
}

func NewHeaderCodec(codec Codec, generator func() Header) *HeaderCodec {
	return new(HeaderCodec).
		setLengthBytes(textLengthEncoding).
		SetCodec(codec).
		SetGenerator(generator)
}

func (c *HeaderCodec) setLengthBytes(lb headerLengthEncoding) *HeaderCodec {
	c.lengthEncoding = lb
	return c
}

// SetCodec 设置头部信息编码
// set the header information code
func (c *HeaderCodec) SetCodec(codec Codec) *HeaderCodec {
	c.codec = codec
	return c
}

// SetGenerator 设置头部结构生成函数
// set the header structure generation function
func (c *HeaderCodec) SetGenerator(generator func() Header) *HeaderCodec {
	c.generator = generator
	return c
}

func (c *HeaderCodec) Encode(writer *bytes.Buffer, h Header) error {
	var p0 [4]byte
	_, _ = writer.Write(p0[:c.lengthEncoding]) // err must be nil

	if h != nil && h.Len() > 0 {
		if err := c.codec.NewEncoder(writer).Encode(h); err != nil {
			return err
		}
	}

	var headerLength = writer.Len() - int(c.lengthEncoding)
	if headerLength > c.lengthEncoding.MaxLength() {
		return errHeaderTooLarge
	}

	var p1 = writer.Bytes()
	var p2 = p1[:c.lengthEncoding]
	if c.lengthEncoding == binaryLengthEncoding {
		binary.BigEndian.PutUint16(p2, uint16(headerLength))
	} else {
		copy(p2, fmt.Sprintf("%04d", headerLength))
	}

	return nil
}

func (c *HeaderCodec) Decode(reader *bytes.Buffer) (Header, error) {
	var v = c.Generate()
	var p [4]byte

	if err := internal.Read(reader, p[:c.lengthEncoding]); err != nil {
		return nil, err
	}

	var headerLength = 0
	if c.lengthEncoding == binaryLengthEncoding {
		headerLength = int(binary.BigEndian.Uint16(p[:c.lengthEncoding]))
	} else {
		n, err := strconv.Atoi(string(p[:c.lengthEncoding]))
		if err != nil {
			return nil, err
		}
		headerLength = n
	}

	if reader.Len() < headerLength {
		return nil, errHeaderSize
	}

	if headerLength > 0 {
		if err := c.codec.Decode(reader.Next(headerLength), v); err != nil {
			return nil, err
		}
	}
	return v.(Header), nil
}

func (c *HeaderCodec) Generate() Header {
	return c.generator()
}

type MapHeader map[string]string

func (c MapHeader) Len() int {
	return len(c)
}

func (c MapHeader) Range(f func(key string, value string)) {
	for k, v := range c {
		f(k, v)
	}
}

func (c MapHeader) Set(key, value string) {
	key = c.formatKey(key)
	c[key] = value
}

func (c MapHeader) Get(key string) string {
	key = c.formatKey(key)
	return c[key]
}

func (c MapHeader) Del(key string) {
	key = c.formatKey(key)
	delete(c, key)
}

func NewHttpHeader(h http.Header) HttpHeader {
	return HttpHeader{Header: h}
}

func (c MapHeader) formatKey(key string) string {
	var b = []byte(key)
	var n = len(b)
	for i := 0; i < n-1; i++ {
		if (i == 0 || b[i-1] == '-') && (b[i] >= 'a' && b[i] <= 'z') {
			b[i] -= 32
		}
	}
	return string(b)
}

type HttpHeader struct {
	http.Header
}

func (c HttpHeader) Len() int {
	return len(c.Header)
}

func (c HttpHeader) Range(f func(key string, value string)) {
	for k, _ := range c.Header {
		f(k, c.Get(k))
	}
}

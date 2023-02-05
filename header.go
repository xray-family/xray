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
	XRealIP     = "X-Real-IP"
	ContentType = "Content-Type"
)

const (
	MimeJson   = "application/json; charset=utf-8"
	MimeStream = "application/octet-stream"
)

const (
	BinaryLengthEncoding HeaderLengthEncoding = 2
	TextLengthEncoding   HeaderLengthEncoding = 4
)

type (
	Header interface {
		Set(key, value string)
		Get(key string) string
		Del(key string)
		Len() int
		Range(f func(key, value string))
	}

	HeaderCodec struct {
		codec          Codec
		lengthEncoding HeaderLengthEncoding
		generator      func() Header
	}

	HeaderLengthEncoding uint8
)

func (c HeaderLengthEncoding) MaxLength() int {
	if c == BinaryLengthEncoding {
		return math.MaxUint16
	}
	return 9999
}

func NewHeaderCodec(lengthEncoding HeaderLengthEncoding, codec Codec, generator func() Header) *HeaderCodec {
	return new(HeaderCodec).
		setLengthBytes(lengthEncoding).
		setCodec(codec).
		SetGenerator(generator)
}

func (c *HeaderCodec) setLengthBytes(lb HeaderLengthEncoding) *HeaderCodec {
	c.lengthEncoding = lb
	return c
}

func (c *HeaderCodec) setCodec(codec Codec) *HeaderCodec {
	c.codec = codec
	return c
}

func (c *HeaderCodec) SetGenerator(generator func() Header) *HeaderCodec {
	c.generator = generator
	return c
}

func (c *HeaderCodec) Encode(writer *bytes.Buffer, h Header) error {
	var p0 [4]byte
	if err := internal.Write(writer, p0[:c.lengthEncoding]); err != nil {
		return err
	}

	if h != nil && h.Len() > 0 {
		if err := c.codec.NewEncoder(writer).Encode(h); err != nil {
			return err
		}
	}

	var headerLength = writer.Len() - int(c.lengthEncoding)
	if headerLength > c.lengthEncoding.MaxLength() {
		return errors.New("encode header error")
	}

	var p1 = writer.Bytes()
	var p2 = p1[:c.lengthEncoding]
	if c.lengthEncoding == BinaryLengthEncoding {
		binary.BigEndian.PutUint16(p2, uint16(headerLength))
	} else {
		copy(p2, fmt.Sprintf("%04d", headerLength))
	}

	return nil
}

func (c *HeaderCodec) Decode(reader internal.BytesReader) (Header, error) {
	var v = c.Generate()
	var p0 [4]byte

	if _, err := reader.Read(p0[:c.lengthEncoding]); err != nil {
		return nil, err
	}

	var headerLength = 0
	if c.lengthEncoding == BinaryLengthEncoding {
		headerLength = int(binary.BigEndian.Uint16(p0[:c.lengthEncoding]))
	} else {
		n, err := strconv.Atoi(string(p0[:c.lengthEncoding]))
		if err != nil {
			return nil, err
		}
		headerLength = n
	}

	var p1 = reader.Bytes()
	if len(p1) < headerLength || headerLength > c.lengthEncoding.MaxLength() {
		return nil, errors.New("decode header error")
	}

	if headerLength > 0 {
		if _, err := reader.Read(p1[:headerLength]); err != nil {
			return nil, err
		}
		if err := c.codec.NewDecoder(bytes.NewReader(p1)).Decode(v); err != nil {
			return nil, err
		}
	}
	return v.(Header), nil
}

func (c *HeaderCodec) Generate() Header {
	return c.generator()
}

type MapHeader map[string]string

func (c MapHeader) Del(key string) {
	delete(c, key)
}

func (c MapHeader) Len() int {
	return len(c)
}

func (c MapHeader) Range(f func(key string, value string)) {
	for k, v := range c {
		f(k, v)
	}
}

func (c MapHeader) Set(key, value string) {
	c[key] = value
}

func (c MapHeader) Get(key string) string {
	return c[key]
}

func NewHttpHeader(h http.Header) HttpHeader {
	return HttpHeader{Header: h}
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

package uRouter

import (
	"bytes"
	"errors"
	"github.com/lxzan/uRouter/constant"
	"github.com/lxzan/uRouter/internal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type headerMocker struct {
	MapHeader
}

func (h headerMocker) Read(p []byte) (n int, err error) {
	return 0, errors.New("test")
}

func (h headerMocker) MarshalJSON() ([]byte, error) {
	return nil, errors.New("test")
}

func TestHeaderCodec(t *testing.T) {
	var as = assert.New(t)

	t.Run("non empty text", func(t *testing.T) {
		var ip = "127.0.0.1"
		var path = "/api/v1"
		var header1 = TextHeader.Generate()
		header1.Set(constant.XRealIP, ip)
		header1.Set(constant.XPath, path)
		var result = bytes.NewBuffer(nil)
		if err := TextHeader.Encode(result, header1); err != nil {
			as.NoError(err)
			return
		}

		header2, err := TextHeader.Decode(result)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal(ip, header2.Get(constant.XRealIP))
		as.Equal(path, header2.Get(constant.XPath))
	})

	t.Run("empty text", func(t *testing.T) {
		var result = bytes.NewBuffer(nil)
		if err := TextHeader.Encode(result, nil); err != nil {
			as.NoError(err)
			return
		}

		header2, err := TextHeader.Decode(result)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal(0, header2.Len())
	})

	t.Run("non empty binary", func(t *testing.T) {
		var ip = "127.0.0.1"
		var path = "/api/v1"
		var header1 = BinaryHeader.Generate()
		header1.Set(constant.XRealIP, ip)
		header1.Set(constant.XPath, path)
		var result = bytes.NewBuffer(nil)
		if err := BinaryHeader.Encode(result, header1); err != nil {
			as.NoError(err)
			return
		}

		header2, err := BinaryHeader.Decode(result)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal(ip, header2.Get(constant.XRealIP))
		as.Equal(path, header2.Get(constant.XPath))
	})

	t.Run("empty binary", func(t *testing.T) {
		var result = bytes.NewBuffer(nil)
		if err := BinaryHeader.Encode(result, nil); err != nil {
			as.NoError(err)
			return
		}

		header2, err := BinaryHeader.Decode(result)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal(0, header2.Len())
	})

	t.Run("encode header error", func(t *testing.T) {
		var header = &headerMocker{MapHeader{}}
		header.Set(constant.ContentType, constant.MimeJson)
		var w = bytes.NewBufferString("")
		as.Error(TextHeader.Encode(w, header))
	})

	t.Run("encode big header", func(t *testing.T) {
		var header = TextHeader.Generate()
		var w = bytes.NewBufferString("")
		for i := 0; i < 1000; i++ {
			var key = string(internal.AlphabetNumeric.Generate(16))
			header.Set(key, "1")
		}
		as.Error(TextHeader.Encode(w, header))
	})

	t.Run("decode small header 1", func(t *testing.T) {
		var buf = bytes.NewBufferString("0012xxx")
		_, err := TextHeader.Decode(buf)
		as.Equal(true, err == errHeaderSize)
	})

	t.Run("decode small header 2", func(t *testing.T) {
		var buf = bytes.NewBufferString("00")
		_, err := TextHeader.Decode(buf)
		as.Error(err)
	})

	t.Run("decode error header 1", func(t *testing.T) {
		var buf = bytes.NewBufferString("00xx")
		_, err := TextHeader.Decode(buf)
		as.Error(err)
	})

	t.Run("decode error header 2", func(t *testing.T) {
		var result = bytes.NewBufferString(`0019{"X-Path":"/greet}{"hello":"你好, 少侠"}`)
		_, err := TextHeader.Decode(result)
		if err != nil {
			as.Error(err)
			return
		}
	})
}

func TestHttpHeader(t *testing.T) {
	var as = assert.New(t)
	var header = HttpHeader{Header: http.Header{}}
	header.Set(constant.ContentType, constant.MimeJson)
	header.Set(constant.XPath, "")
	as.Equal(2, header.Len())

	var keys []string
	header.Del(constant.XPath)
	header.Range(func(key string, value string) {
		keys = append(keys, key)
	})
	as.Equal(1, len(keys))
}

func TestMapHeader(t *testing.T) {
	var as = assert.New(t)
	var header = MapHeader{}
	header.Set(constant.ContentType, constant.MimeJson)
	header.Set(constant.XPath, "")
	as.Equal(2, header.Len())

	var keys []string
	header.Del(constant.XPath)
	header.Range(func(key string, value string) {
		keys = append(keys, key)
	})
	as.Equal(1, len(keys))

	as.Equal(constant.XPath, header.formatKey("x-path"))
	as.Equal(constant.ContentType, header.formatKey("content-type"))
	as.Equal("Ct", header.formatKey("ct"))
	as.Equal("", header.formatKey(""))
}

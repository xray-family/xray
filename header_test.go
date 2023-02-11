package uRouter

import (
	"bytes"
	"errors"
	"github.com/lxzan/uRouter/internal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type headerMocker struct {
	*MapHeader
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
		var header1 = TextMapHeader.Generate()
		header1.Set(internal.XRealIP, ip)
		header1.Set(UPath, path)
		var result = bytes.NewBuffer(nil)
		if err := TextMapHeader.Encode(result, header1); err != nil {
			as.NoError(err)
			return
		}

		header2, err := TextMapHeader.Decode(result)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal(ip, header2.Get(internal.XRealIP))
		as.Equal(path, header2.Get(UPath))
	})

	t.Run("empty text", func(t *testing.T) {
		var result = bytes.NewBuffer(nil)
		if err := TextMapHeader.Encode(result, nil); err != nil {
			as.NoError(err)
			return
		}

		header2, err := TextMapHeader.Decode(result)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal(0, header2.Len())
	})

	t.Run("non empty binary", func(t *testing.T) {
		var ip = "127.0.0.1"
		var path = "/api/v1"
		var header1 = BinaryMapHeader.Generate()
		header1.Set(internal.XRealIP, ip)
		header1.Set(UPath, path)
		var result = bytes.NewBuffer(nil)
		if err := BinaryMapHeader.Encode(result, header1); err != nil {
			as.NoError(err)
			return
		}

		header2, err := BinaryMapHeader.Decode(result)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal(ip, header2.Get(internal.XRealIP))
		as.Equal(path, header2.Get(UPath))
	})

	t.Run("empty binary", func(t *testing.T) {
		var result = bytes.NewBuffer(nil)
		if err := BinaryMapHeader.Encode(result, nil); err != nil {
			as.NoError(err)
			return
		}

		header2, err := BinaryMapHeader.Decode(result)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal(0, header2.Len())
	})

	t.Run("encode header error", func(t *testing.T) {
		var header = &headerMocker{NewMapHeader()}
		header.Set(internal.ContentType, internal.MimeJson)
		var w = bytes.NewBufferString("")
		as.Error(TextMapHeader.Encode(w, header))
	})

	t.Run("encode big header", func(t *testing.T) {
		var header = TextMapHeader.Generate()
		var w = bytes.NewBufferString("")
		for i := 0; i < 1000; i++ {
			var key = string(internal.AlphabetNumeric.Generate(16))
			header.Set(key, "1")
		}
		as.Error(TextMapHeader.Encode(w, header))
	})

	t.Run("decode small header 1", func(t *testing.T) {
		var buf = bytes.NewBufferString("0012xxx")
		_, err := TextMapHeader.Decode(buf)
		as.Equal(true, err == internal.ErrParseHeader)
	})

	t.Run("decode small header 2", func(t *testing.T) {
		var buf = bytes.NewBufferString("00")
		_, err := TextMapHeader.Decode(buf)
		as.Error(err)
	})

	t.Run("decode error header 1", func(t *testing.T) {
		var buf = bytes.NewBufferString("00xx")
		_, err := TextMapHeader.Decode(buf)
		as.Error(err)
	})

	t.Run("decode error header 2", func(t *testing.T) {
		var result = bytes.NewBufferString(`0019{"X-Path":"/greet}{"hello":"你好, 少侠"}`)
		_, err := TextMapHeader.Decode(result)
		if err != nil {
			as.Error(err)
			return
		}
	})
}

func TestHttpHeader(t *testing.T) {
	var as = assert.New(t)
	var header = HttpHeader{Header: http.Header{}}
	header.Set(internal.ContentType, internal.MimeJson)
	header.Set(UPath, "")
	as.Equal(2, header.Len())

	var keys []string
	header.Del(UPath)
	header.Range(func(key string, value string) {
		keys = append(keys, key)
	})
	as.Equal(1, len(keys))
}

func TestMapHeader(t *testing.T) {
	var as = assert.New(t)
	var header = mapHeaderTemplate.Generate().(*MapHeader)
	header.Set(internal.ContentType, internal.MimeJson)
	header.Set(UPath, "")
	as.Equal(2, header.Len())

	var keys []string
	header.Del(UPath)
	header.Range(func(key string, value string) {
		keys = append(keys, key)
	})
	as.Equal(1, len(keys))

	as.Equal(UPath, header.formatKey("u-path"))
	as.Equal(internal.ContentType, header.formatKey("content-type"))
	as.Equal("Ct", header.formatKey("ct"))
	as.Equal("", header.formatKey(""))
}

func TestHeader_Reset(t *testing.T) {

	t.Run("map header", func(t *testing.T) {
		var m = NewMapHeader()
		for i := 0; i < 64; i++ {
			key := internal.AlphabetNumeric.Generate(16)
			m.Set(string(key), "1")
		}
		m.Reset()
		assert.Equal(t, 0, m.Len())
		m.Close()
	})

	t.Run("http header", func(t *testing.T) {
		var m = &HttpHeader{Header: http.Header{}}
		for i := 0; i < 64; i++ {
			key := internal.AlphabetNumeric.Generate(16)
			m.Set(string(key), "1")
		}
		m.Reset()
		assert.Equal(t, 0, m.Len())
		m.Close()
		m.Generate()
	})

}

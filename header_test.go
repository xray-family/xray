package uRouter

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHeaderCodec(t *testing.T) {
	var as = assert.New(t)

	t.Run("non empty text", func(t *testing.T) {
		var ip = "127.0.0.1"
		var path = "/api/v1"
		var header1 = TextHeader.Generate()
		header1.Set(XRealIP, ip)
		header1.Set(XPath, path)
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
		as.Equal(ip, header2.Get(XRealIP))
		as.Equal(path, header2.Get(XPath))
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
		header1.Set(XRealIP, ip)
		header1.Set(XPath, path)
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
		as.Equal(ip, header2.Get(XRealIP))
		as.Equal(path, header2.Get(XPath))
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
}

func TestHttpHeader(t *testing.T) {
	var as = assert.New(t)
	var header = HttpHeader{Header: http.Header{}}
	header.Set(ContentType, MimeJson)
	header.Set(XPath, "")
	as.Equal(2, header.Len())

	var keys []string
	header.Del(XPath)
	header.Range(func(key string, value string) {
		keys = append(keys, key)
	})
	as.Equal(1, len(keys))
}

func TestMapHeader(t *testing.T) {
	var as = assert.New(t)
	var header = MapHeader{}
	header.Set(ContentType, MimeJson)
	header.Set(XPath, "")
	as.Equal(2, header.Len())

	var keys []string
	header.Del(XPath)
	header.Range(func(key string, value string) {
		keys = append(keys, key)
	})
	as.Equal(1, len(keys))
}

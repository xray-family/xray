package uRouter

import (
	"bytes"
	"github.com/stretchr/testify/assert"
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
		as.Equal(0, len(header2.(F)))
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
		as.Equal(0, len(header2.(F)))
	})
}

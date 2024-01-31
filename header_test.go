package xray

import (
	"github.com/stretchr/testify/assert"
	"github.com/xray-family/xray/internal"
	"net/http"
	"testing"
)

func TestHttpHeader(t *testing.T) {
	var as = assert.New(t)
	var header = HttpHeader{Header: http.Header{}}
	header.Set(ContentType, MimeJson)
	header.Set(XPath, "")
	as.Equal(2, header.Len())

	var keys []string
	header.Del(XPath)
	header.Range(func(key string, value string) bool {
		keys = append(keys, key)
		return true
	})
	as.Equal(1, len(keys))
	as.Nil(header.New())
}

func TestSliceHeader(t *testing.T) {
	tpl := &SliceHeader{}

	t.Run("del", func(t *testing.T) {
		h := tpl.New()
		h.Set("x-token", "1")
		h.Set("x-path", "2")
		assert.Equal(t, h.Len(), 2)
		h.Del("x-token")
		h.Del("x-token")
		h.Del("x-path")
		assert.Equal(t, h.Len(), 0)
		assert.Equal(t, h.Get("x-path"), "")
	})

	t.Run("set", func(t *testing.T) {
		h := tpl.New()
		h.Set("x-token", "1")
		h.Set("x-token", "2")
		h.Set("x-real-ip", "3")
		assert.Equal(t, h.Get("x-token"), "2")
		assert.Equal(t, h.Len(), 2)
	})

	t.Run("range", func(t *testing.T) {
		h := tpl.New()
		a := make([]string, 0)
		h.Set("x-token", "1")
		h.Set("x-path", "2")
		h.Set("x-real-ip", "3")
		h.Range(func(key, value string) bool {
			a = append(a, value)
			return len(a) < 2
		})
		assert.ElementsMatch(t, a, []string{"1", "2"})
	})
}

func TestSliceHeader_Random(t *testing.T) {
	var sh = &SliceHeader{}
	var hashmap = make(map[string]string)
	for i := 0; i < 10000; i++ {
		var flag = internal.Numeric.Intn(10)
		var key = string(internal.Numeric.Generate(3))
		var val = string(internal.Numeric.Generate(6))
		switch flag {
		case 0, 1, 2:
			sh.Set(key, val)
			hashmap[key] = val
		case 3:
			sh.Del(key)
			delete(hashmap, key)
		}
	}
	assert.Equal(t, sh.Len(), len(hashmap))
	sh.Range(func(key string, value string) bool {
		assert.Equal(t, value, hashmap[key])
		return true
	})
}

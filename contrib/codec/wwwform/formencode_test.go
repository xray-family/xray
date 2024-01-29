package wwwform

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCodec_Decode(t *testing.T) {
	as := assert.New(t)

	type Message struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}

	t.Run("decoder", func(t *testing.T) {
		var s = "name=xxx&age=2"
		var v = Message{}
		err := Codec.NewDecoder(strings.NewReader(s)).Decode(&v)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal("xxx", v.Name)
		as.Equal(2, v.Age)
	})

	t.Run("decode", func(t *testing.T) {
		var s = "name=xxx&age=2"
		var v = Message{}
		err := Codec.Decode([]byte(s), &v)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal("xxx", v.Name)
		as.Equal(2, v.Age)
	})

	t.Run("decode from string", func(t *testing.T) {
		var s = "name=xxx&age=2"
		var v = Message{}
		err := Codec.DecodeFromString(s, &v)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal("xxx", v.Name)
		as.Equal(2, v.Age)
	})
}

func TestCodec_Encode(t *testing.T) {
	as := assert.New(t)

	type Message struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}

	t.Run("encoder", func(t *testing.T) {
		var v = Message{
			Name: "小明",
			Age:  12,
		}
		var buf = bytes.NewBufferString("")
		err := Codec.NewEncoder(buf).Encode(&v)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal("age=12&name=%E5%B0%8F%E6%98%8E", buf.String())
	})

	t.Run("encode", func(t *testing.T) {
		var v = Message{
			Name: "小明",
			Age:  12,
		}
		b, err := Codec.Encode(v)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal("age=12&name=%E5%B0%8F%E6%98%8E", string(b))
	})

	t.Run("encode to string", func(t *testing.T) {
		var v = Message{
			Name: "小明",
			Age:  12,
		}
		b, err := Codec.EncodeToString(v)
		if err != nil {
			as.NoError(err)
			return
		}
		as.Equal("age=12&name=%E5%B0%8F%E6%98%8E", b)
	})
}

package xray

import (
	"bytes"
	"errors"
	"github.com/lxzan/xray/codec"
	"github.com/lxzan/xray/constant"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func newContextMocker(options ...Option) *Context {
	var router = New(options...)
	var request = &Request{
		Header: HttpHeader{Header: http.Header{}},
		Body:   bytes.NewBuffer(nil),
	}
	var writer = newResponseWriterMocker()
	var ctx = NewContext(router, request, writer)
	return ctx
}

func newResponseWriterMocker() *responseWriterMocker {
	return &responseWriterMocker{
		protocol:   constant.ProtocolHTTP,
		statusCode: 0,
		header:     HttpHeader{Header: http.Header{}},
		buf:        bytes.NewBuffer(nil),
	}
}

type responseWriterMocker struct {
	protocol   string
	statusCode int
	header     HttpHeader
	buf        *bytes.Buffer
}

func (c *responseWriterMocker) SetProtocol(p string) {
	c.protocol = p
}

func (c *responseWriterMocker) Protocol() string {
	return c.protocol
}

func (c *responseWriterMocker) Header() Header {
	return c.header
}

func (c *responseWriterMocker) Write(p []byte) (int, error) {
	if strings.TrimSpace(string(p)) == `"err"` {
		return 0, errors.New("test error")
	}
	return c.buf.Write(p)
}

func (c *responseWriterMocker) Code(code int) {
	c.statusCode = code
}

func (c *responseWriterMocker) Flush() error {
	return nil
}

func (c *responseWriterMocker) Raw() any {
	return nil
}

func TestContext_BindJSON(t *testing.T) {
	var as = assert.New(t)

	t.Run("", func(t *testing.T) {
		var ctx = newContextMocker()
		ctx.Request.Body = bytes.NewBufferString(`{"age":1}`)
		var params = struct {
			Age int `json:"age"`
		}{}
		as.NoError(ctx.BindJSON(&params))
		as.Equal(1, params.Age)
	})

	t.Run("", func(t *testing.T) {
		var ctx = newContextMocker()
		ctx.Request.Body = strings.NewReader(`{"age":1}`)
		var params = struct {
			Age int `json:"age"`
		}{}
		as.NoError(ctx.BindJSON(&params))
		as.Equal(1, params.Age)
	})

	t.Run("", func(t *testing.T) {
		var ctx = newContextMocker()
		ctx.Request.Body = bytes.NewBufferString(`{"age":"1}`)
		var params = struct {
			Age int `json:"age"`
		}{}
		as.Error(ctx.BindJSON(&params))
	})

	t.Run("", func(t *testing.T) {
		var ctx = newContextMocker()
		var params = struct {
			Age int `json:"age"`
		}{}
		as.Error(ctx.BindJSON(&params))
	})
}

func TestContext_Write(t *testing.T) {
	var as = assert.New(t)

	t.Run("write json 1", func(t *testing.T) {
		var ctx = newContextMocker()
		var params = Any{"name": "aha"}
		if err := ctx.WriteJSON(http.StatusOK, params); err != nil {
			as.NoError(err)
			return
		}
		var writer = ctx.Writer.(*responseWriterMocker)
		as.Equal(http.StatusOK, writer.statusCode)
		as.Equal(constant.MimeJson, writer.header.Get(constant.ContentType))
		var buf = bytes.NewBufferString("")
		codec.StdJsonCodec.NewEncoder(buf).Encode(params)
		as.Equal(buf.Len(), writer.buf.Len())
	})

	t.Run("write string", func(t *testing.T) {
		var ctx = newContextMocker()
		var params = "hello"
		if err := ctx.WriteString(http.StatusOK, params); err != nil {
			as.NoError(err)
			return
		}
		var writer = ctx.Writer.(*responseWriterMocker)
		as.Equal(http.StatusOK, writer.statusCode)
		as.Equal("", writer.header.Get(constant.ContentType))
		as.Equal(params, writer.buf.String())
	})

	t.Run("write string", func(t *testing.T) {
		var ctx = newContextMocker()
		var params = []byte("hello")
		if err := ctx.WriteBytes(http.StatusOK, params); err != nil {
			as.NoError(err)
			return
		}
		var writer = ctx.Writer.(*responseWriterMocker)
		as.Equal(http.StatusOK, writer.statusCode)
		as.Equal("", writer.header.Get(constant.ContentType))
		as.Equal(string(params), writer.buf.String())
	})

	t.Run("write bytes", func(t *testing.T) {
		var ctx = newContextMocker()
		as.Error(ctx.WriteBytes(http.StatusOK, []byte(`"err"`)))
	})

	t.Run("write json", func(t *testing.T) {
		var ctx = newContextMocker()
		var err = ctx.WriteJSON(200, "err")
		as.Error(err)
	})
}

func TestContext_Storage(t *testing.T) {
	var as = assert.New(t)
	var ctx = newContextMocker()
	ctx.Set("name", "aha")
	ctx.Set("age", 1)
	{
		v, _ := ctx.Get("name")
		as.Equal("aha", v)
	}
	{
		v, _ := ctx.Get("age")
		as.Equal(1, v)
	}
}

func TestContext_Param(t *testing.T) {
	var as = assert.New(t)
	var router = New()

	t.Run("", func(t *testing.T) {
		var ctx = NewContext(router, &Request{
			Header: &HttpHeader{http.Header{}},
			VPath:  "/:id",
		}, newResponseWriterMocker())
		id := ctx.Param("id")
		as.Equal("", id)
	})

	t.Run("", func(t *testing.T) {
		var ctx = NewContext(router, &Request{
			Header: &HttpHeader{http.Header{}},
			VPath:  "/api/v1",
		}, newResponseWriterMocker())
		id := ctx.Param("id")
		as.Equal("", id)
	})

	t.Run("", func(t *testing.T) {
		var ctx = NewContext(router, &Request{
			VPath: "/api/v1",
			RPath: "/api/v1",
		}, newResponseWriterMocker())
		id := ctx.Param("id")
		as.Equal("", id)
	})
}

func TestRequest_Close(t *testing.T) {
	var r = &Request{Header: &HttpHeader{http.Header{}}, Body: bytes.NewBufferString("")}
	r.Close()
	assert.Nil(t, r.Body)
}

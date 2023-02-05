package uRouter

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func newContextMocker() *Context {
	var request = &Request{
		Header: HttpHeader{Header: http.Header{}},
		Body:   bytes.NewBuffer(nil),
	}
	var writer = newResponseWriterMocker()
	var ctx = NewContext(request, writer)
	return ctx
}

func newResponseWriterMocker() ResponseWriter {
	return &responseWriterMocker{
		statusCode: 0,
		header:     HttpHeader{Header: http.Header{}},
		buf:        bytes.NewBuffer(nil),
	}
}

type responseWriterMocker struct {
	statusCode int
	header     HttpHeader
	buf        *bytes.Buffer
}

func (c *responseWriterMocker) Protocol() string {
	return ProtocolHTTP
}

func (c *responseWriterMocker) Header() Header {
	return c.header
}

func (c *responseWriterMocker) Write(p []byte) (int, error) {
	return c.buf.Write(p)
}

func (c *responseWriterMocker) Code(code int) {
	c.statusCode = code
}

func (c *responseWriterMocker) Flush() error {
	return nil
}

func (c *responseWriterMocker) Raw() interface{} {
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
		ctx.Request.Body = bytes.NewBufferString(`{"age":"1}`)
		var params = struct {
			Age int `json:"age"`
		}{}
		as.Error(ctx.BindJSON(&params))
	})

	t.Run("", func(t *testing.T) {
		var ctx = newContextMocker()
		ctx.Request.Body = nil
		var params = struct {
			Age int `json:"age"`
		}{}
		as.Error(ctx.BindJSON(&params))
	})
}

func TestContext_Write(t *testing.T) {
	var as = assert.New(t)

	t.Run("write json", func(t *testing.T) {
		var ctx = newContextMocker()
		var params = Any{"name": "aha"}
		if err := ctx.WriteJSON(http.StatusOK, params); err != nil {
			as.NoError(err)
			return
		}
		var writer = ctx.Writer.(*responseWriterMocker)
		as.Equal(http.StatusOK, writer.statusCode)
		as.Equal(MimeJson, writer.header.Get(ContentType))
		var buf = bytes.NewBufferString("")
		defaultJsonCodec.NewEncoder(buf).Encode(params)
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
		as.Equal("", writer.header.Get(ContentType))
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
		as.Equal("", writer.header.Get(ContentType))
		as.Equal(string(params), writer.buf.String())
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

func TestContext_Others(t *testing.T) {
	var as = assert.New(t)
	var ctx = newContextMocker()
	as.Nil(ctx.Request.Raw)
	as.Nil(ctx.Writer.Raw())
}

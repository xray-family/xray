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
		protocol:   ProtocolHTTP,
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

	t.Run("write json 1", func(t *testing.T) {
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

	t.Run("write json 2", func(t *testing.T) {
		var ctx = newContextMocker()
		var header = &headerMocker{MapHeader{}}
		header.Set(ContentType, MimeJson)
		as.Error(ctx.WriteJSON(http.StatusOK, header))
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

	t.Run("write reader", func(t *testing.T) {
		var ctx = newContextMocker()
		var header = &headerMocker{MapHeader{}}
		header.Set(ContentType, MimeJson)
		as.Error(ctx.WriteReader(http.StatusOK, header))
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
	SetJsonCodec(StdJsonCodec)
	SetBufferPool(newBufferPool())
	defaultGenerator()
	as.Nil(ctx.Request.Raw)
	as.Nil(ctx.Writer.Raw())
}

func TestContext_Param(t *testing.T) {
	var as = assert.New(t)

	t.Run("", func(t *testing.T) {
		var ctx = NewContext(&Request{
			Header: NewHttpHeader(http.Header{
				XPath: []string{""},
			}),
			VPath: "/:id",
		}, newResponseWriterMocker())
		id := ctx.Param("id")
		as.Equal("", id)
	})

	t.Run("", func(t *testing.T) {
		var ctx = NewContext(&Request{
			Header: NewHttpHeader(http.Header{
				XPath: []string{"/api/v1"},
			}),
			VPath: "/api/v1",
		}, newResponseWriterMocker())
		id := ctx.Param("id")
		as.Equal("", id)
	})
}

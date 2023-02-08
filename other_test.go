package uRouter

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestBufferPool(t *testing.T) {
	var as = assert.New(t)
	var p = DefaultBufferPool()
	SetBufferPool(newBufferPool())
	p.SetSize(4 * 1024)
	as.NotNil(p.Get())
	p.Put(bytes.NewBufferString(""))
	as.NotNil(p.Get())
}

func TestAccessLog(t *testing.T) {
	var r = New()
	r.On("test", AccessLog())
	var ctx = newContextMocker()
	ctx.Request.Header.Set(XPath, "test")
	r.Emit(ctx)
}

func TestWebSocket(t *testing.T) {
	t.Run("reject", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.Use(WebSocket())
		r.On("test", func(ctx *Context) {
			sum++
		})
		var ctx = newContextMocker()
		ctx.Request.Header.Set(XPath, "test")
		r.Emit(ctx)
		assert.Equal(t, 0, sum)
	})

	t.Run("pass", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.Use(WebSocket())
		r.On("test", func(ctx *Context) {
			sum++
		})
		var ctx = newContextMocker()
		ctx.Writer.(*responseWriterMocker).SetProtocol(ProtocolWebSocket)
		ctx.Request.Header.Set(XPath, "test")
		r.Emit(ctx)
		assert.Equal(t, 1, sum)
	})
}

func TestHTTP(t *testing.T) {
	t.Run("reject 1", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.Use(HTTP(http.MethodPost))
		r.On("test", func(ctx *Context) {
			sum++
		})

		var ctx = newContextMocker()
		ctx.Request.Raw = &http.Request{Method: http.MethodGet}
		ctx.Request.Header.Set(XPath, "test")
		r.Emit(ctx)
		assert.Equal(t, 0, sum)
	})

	t.Run("reject 2", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.Use(HTTP(http.MethodPost))
		r.On("test", func(ctx *Context) {
			sum++
		})

		var ctx = newContextMocker()
		ctx.Writer.(*responseWriterMocker).SetProtocol(ProtocolWebSocket)
		ctx.Request.Raw = &http.Request{Method: http.MethodPost}
		ctx.Request.Header.Set(XPath, "test")
		r.Emit(ctx)
		assert.Equal(t, 0, sum)
	})

	t.Run("pass", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.Use(HTTP(http.MethodPost))
		r.On("test", func(ctx *Context) {
			sum++
		})

		var ctx = newContextMocker()
		ctx.Request.Raw = &http.Request{Method: http.MethodPost}
		ctx.Request.Header.Set(XPath, "test")
		r.Emit(ctx)
		assert.Equal(t, 1, sum)
	})
}

func TestRecovery(t *testing.T) {
	var as = assert.New(t)

	t.Run("with recovery", func(t *testing.T) {
		var r = New()

		defer func() {
			e := recover()
			as.Nil(e)
		}()

		r.Use(Recovery())
		r.On("test", func(ctx *Context) {
			panic("1")
		})
		var ctx = newContextMocker()
		ctx.Request.Header.Set(XPath, "test")
		r.Emit(ctx)
	})

	t.Run("without recovery", func(t *testing.T) {
		var r = New()

		defer func() {
			e := recover()
			as.NotNil(e)
		}()

		r.On("test", func(ctx *Context) {
			panic("1")
		})
		var ctx = newContextMocker()
		ctx.Request.Header.Set(XPath, "test")
		r.Emit(ctx)
	})
}

type closerMocker struct{}

func (c *closerMocker) Close() {}

func TestClose(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var c = io.NopCloser(bytes.NewBuffer(nil))
		Close(c)
	})

	t.Run("", func(t *testing.T) {
		var c = &closerMocker{}
		Close(c)
	})

	t.Run("", func(t *testing.T) {
		var c = bytes.NewBuffer(nil)
		Close(c)
	})
}

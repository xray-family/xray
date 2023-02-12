package uRouter

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestAccessLog(t *testing.T) {
	var r = New()
	r.On("test", AccessLog())
	var ctx = newContextMocker()
	r.Emit("test", ctx)
}

func TestWebSocket(t *testing.T) {
	t.Run("reject", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.Use(WebSocketRequired())
		r.On("test", func(ctx *Context) {
			sum++
		})
		var ctx = newContextMocker()
		ctx.Request.Header.Set(UPath, "test")
		r.Emit("test", ctx)
		assert.Equal(t, 0, sum)
	})

	t.Run("pass", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.Use(WebSocketRequired())
		r.OnEvent(http.MethodGet, "/test", func(ctx *Context) {
			sum++
		}, AccessLog())
		r.OnEvent(http.MethodPost, "/test", func(ctx *Context) {
			sum += 2
		})
		r.OnEvent(http.MethodPost, "/aha", func(ctx *Context) {
			sum += 4
		})
		r.Start()

		var ctx = newContextMocker()
		ctx.Writer.(*responseWriterMocker).SetProtocol(ProtocolWebSocket)
		ctx.Request.Header.Set(UPath, "/test")
		r.EmitEvent(http.MethodGet, "/test", ctx)
		assert.Equal(t, 1, sum)
	})

	t.Run("recovery", func(t *testing.T) {
		defer func() {
			e := recover()
			assert.Nil(t, e)
		}()

		r := New()
		r.Use(Recovery())
		var path = "/test"
		r.On(path, func(ctx *Context) {
			panic("recovery test")
		})
		r.StartSilently()

		r.Emit(path, newContextMocker())
	})
}

func TestHTTP(t *testing.T) {
	t.Run("reject 1", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.Use(HttpRequired(http.MethodPost))
		r.On("test", func(ctx *Context) {
			sum++
		})

		var ctx = newContextMocker()
		ctx.Request.Raw = &http.Request{Method: http.MethodGet}
		ctx.Request.Header.Set(UPath, "test")
		r.Emit("test", ctx)
		assert.Equal(t, 0, sum)
	})

	t.Run("reject 2", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.Use(HttpRequired(http.MethodPost))
		r.On("test", func(ctx *Context) {
			sum++
		})

		var ctx = newContextMocker()
		ctx.Writer.(*responseWriterMocker).SetProtocol(ProtocolWebSocket)
		ctx.Request.Raw = &http.Request{Method: http.MethodPost}
		ctx.Request.Header.Set(UPath, "test")
		r.Emit("test", ctx)
		assert.Equal(t, 0, sum)
	})

	t.Run("reject 2", func(t *testing.T) {
		var r = New()
		var sum = 0
		var path = "/test"
		r.Use(HttpRequired(http.MethodPost))
		r.OnEvent(http.MethodGet, path, func(ctx *Context) {
			sum++
		})
		r.StartSilently()

		var ctx = newContextMocker()
		ctx.Request.Raw = &http.Request{}
		r.EmitEvent(http.MethodPost, path, ctx)
		assert.Equal(t, 0, sum)
	})

	t.Run("reject 3", func(t *testing.T) {
		var r = New()
		var sum = 0
		var path = "/test"
		r.Use(HttpRequired(http.MethodPost))
		r.OnEvent(http.MethodGet, path, func(ctx *Context) {
			sum++
		})
		r.StartSilently()

		var ctx = newContextMocker()
		ctx.Request.Raw = &http.Request{}
		w := newResponseWriterMocker()
		w.SetProtocol(ProtocolWebSocket)
		ctx.Writer = w
		r.EmitEvent(http.MethodGet, path, ctx)
		assert.Equal(t, 0, sum)
	})

	t.Run("pass", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.Use(HttpRequired(http.MethodPost))
		r.On("/test", func(ctx *Context) { sum++ })
		r.StartSilently()

		var ctx = newContextMocker()
		ctx.Request.Raw = &http.Request{Method: http.MethodPost}
		ctx.Request.Header.Set(UPath, "/test")
		r.Emit("/test", ctx)
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
		r.On("/test", func(ctx *Context) {
			panic("1")
		})
		var ctx = newContextMocker()
		ctx.Request.Header.Set(UPath, "/test")
		r.Emit("/test", ctx)
	})

	t.Run("without recovery", func(t *testing.T) {
		var r = New()

		defer func() {
			e := recover()
			as.NotNil(e)
		}()

		r.On("/test", func(ctx *Context) {
			panic("1")
		})
		r.StartSilently()

		var ctx = newContextMocker()
		ctx.Request.Header.Set(UPath, "/test")
		r.Emit("/test", ctx)
		println(1)
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

func TestLogger(t *testing.T) {
	defer func() {
		recover()
	}()

	SetLogger(defaultLogger)
	Logger().Debug("1")
	Logger().Info("2")
	Logger().Warn("3")
	Logger().Error("4")
	Logger().Panic("5")
}

func TestAny(t *testing.T) {
	var as = assert.New(t)
	var m = Any{
		"num1": 1,
		"num2": int64(2),
		"str":  "3",
	}
	as.Equal(1, m.ToInt("num1"))
	as.Equal(int64(2), m.ToInt64("num2"))
	as.Equal("3", m.ToString("str"))

	as.Equal(false, m.Exists("xxx"))
	as.Equal(0, m.ToInt("xxx"))
	as.Equal(int64(0), m.ToInt64("xxx"))
	as.Equal("", m.ToString("xxx"))
}

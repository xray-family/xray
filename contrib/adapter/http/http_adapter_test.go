package http

import (
	"bytes"
	"github.com/lxzan/uRouter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func newWriterMocker() http.ResponseWriter {
	return &writerMocker{
		header: http.Header{},
		buf:    bytes.NewBufferString(""),
	}
}

type writerMocker struct {
	header http.Header
	buf    *bytes.Buffer
	code   int
}

func (c *writerMocker) Header() http.Header {
	return c.header
}

func (c *writerMocker) Write(p []byte) (int, error) {
	return c.buf.Write(p)
}

func (c *writerMocker) WriteHeader(statusCode int) {
	c.code = statusCode
}

func TestNewAdapter(t *testing.T) {
	var as = assert.New(t)

	t.Run("abort", func(t *testing.T) {
		var sum = int64(0)
		var router = uRouter.New()

		router.Use(func(ctx *uRouter.Context) {
			return
		})

		router.On("/test", func(ctx *uRouter.Context) {
			sum++
		})

		var adapter = NewAdapter(router)
		adapter.ServeHTTP(nil, &http.Request{
			Header: http.Header{},
			URL: &url.URL{
				Path: "/test",
			}})
		as.Equal(int64(0), sum)
	})

	t.Run("next", func(t *testing.T) {
		var sum = int64(0)
		var router = uRouter.New()

		router.Use(func(ctx *uRouter.Context) {
			ctx.Next()
			return
		})

		router.On("/test", func(ctx *uRouter.Context) {
			sum++
		})

		var adapter = NewAdapter(router)
		adapter.ServeHTTP(nil, &http.Request{
			Header: http.Header{},
			URL: &url.URL{
				Path: "/test",
			}})

		as.Equal(int64(1), sum)
	})

	t.Run("complex", func(t *testing.T) {
		var router = uRouter.New()
		var adapter = NewAdapter(router)

		router.Use(func(ctx *uRouter.Context) {
			ctx.Set("sum", 0)
			ctx.Next()
		})

		router.OnNoMatch = func(ctx *uRouter.Context) {
			v, _ := ctx.Get("sum")
			as.Equal(0, v.(int))
		}

		g0 := router.Group("api/v1", func(ctx *uRouter.Context) {
			v, _ := ctx.Get("sum")
			ctx.Set("sum", v.(int)+1)
			ctx.Next()
		})

		g1 := g0.Group("user", func(ctx *uRouter.Context) {
			v, _ := ctx.Get("sum")
			ctx.Set("sum", v.(int)+4)
			ctx.Next()
		})

		g0.On("/t1", func(ctx *uRouter.Context) {
			v, _ := ctx.Get("sum")
			as.Equal(3, v.(int))

			{
				ctx.Writer.Header().Set(uRouter.ContentType, "plain/text")
				as.NoError(ctx.WriteString(http.StatusOK, "OK"))
				_, ok := ctx.RawResponseWriter().(http.ResponseWriter)
				as.Equal(true, ok)
				as.Equal("plain/text", ctx.Writer.Header().Get(uRouter.ContentType))
			}

		}, func(ctx *uRouter.Context) {
			v, _ := ctx.Get("sum")
			ctx.Set("sum", v.(int)+2)
			ctx.Next()
		})

		g0.On("t2", func(ctx *uRouter.Context) {
			v, _ := ctx.Get("sum")
			as.Equal(1, v.(int))
		})

		g1.On("t3", func(ctx *uRouter.Context) {
			v, _ := ctx.Get("sum")
			as.Equal(5, v.(int))
		})

		g2 := g0.Group("session")

		g2.On("t4", func(ctx *uRouter.Context) {
			v, _ := ctx.Get("sum")
			as.Equal(1, v.(int))
		})

		router.Display()

		adapter.ServeHTTP(newWriterMocker(), &http.Request{Header: http.Header{}, URL: &url.URL{Path: "/0123abc"}})
		adapter.ServeHTTP(newWriterMocker(), &http.Request{Header: http.Header{}, URL: &url.URL{Path: "/api/v1/t1"}})
		adapter.ServeHTTP(newWriterMocker(), &http.Request{Header: http.Header{}, URL: &url.URL{Path: "/api/v1/t2"}})
		adapter.ServeHTTP(newWriterMocker(), &http.Request{Header: http.Header{}, URL: &url.URL{Path: "/api/v1/user/t3"}})
		adapter.ServeHTTP(newWriterMocker(), &http.Request{Header: http.Header{}, URL: &url.URL{Path: "/api/v1/session/t4"}})

		adapter.ServeHTTP(newWriterMocker(), &http.Request{Header: http.Header{}, URL: &url.URL{Path: "/0123abc"}})
		adapter.ServeHTTP(newWriterMocker(), &http.Request{Header: http.Header{}, URL: &url.URL{Path: "/api/v1/t1"}})
		adapter.ServeHTTP(newWriterMocker(), &http.Request{Header: http.Header{}, URL: &url.URL{Path: "/api/v1/t2"}})
		adapter.ServeHTTP(newWriterMocker(), &http.Request{Header: http.Header{}, URL: &url.URL{Path: "/api/v1/user/t3"}})
		adapter.ServeHTTP(newWriterMocker(), &http.Request{Header: http.Header{}, URL: &url.URL{Path: "/api/v1/session/t4"}})
	})
}

package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/xray-family/xray"
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
		var router = xray.New()
		var adapter = NewAdapter(router)

		router.Use(func(ctx *xray.Context) {
			return
		})

		router.On("/test", func(ctx *xray.Context) {
			sum++
		})

		adapter.ServeHTTP(nil, &http.Request{
			Header: http.Header{},
			URL: &url.URL{
				Path: "/test",
			}})
		as.Equal(int64(0), sum)
	})

	t.Run("next", func(t *testing.T) {
		var sum = int64(0)
		var router = xray.New()
		var adapter = NewAdapter(router)

		router.Use(func(ctx *xray.Context) {
			ctx.Next()
			return
		})

		router.On("/test", func(ctx *xray.Context) {
			sum++
		})

		adapter.ServeHTTP(nil, &http.Request{
			Header: http.Header{},
			URL: &url.URL{
				Path: "/test",
			}})

		as.Equal(int64(1), sum)
	})

	t.Run("complex", func(t *testing.T) {
		var router = xray.New()
		var adapter = NewAdapter(router)

		router.Use(func(ctx *xray.Context) {
			ctx.Set("sum", 0)
			ctx.Next()
		})

		router.SetHandlerNotFound(func(ctx *xray.Context) {
			v, _ := ctx.Get("sum")
			as.Equal(0, v.(int))
		})

		g0 := router.Group("api/v1", func(ctx *xray.Context) {
			v, _ := ctx.Get("sum")
			ctx.Set("sum", v.(int)+1)
			ctx.Next()
		})

		g1 := g0.Group("user", func(ctx *xray.Context) {
			v, _ := ctx.Get("sum")
			ctx.Set("sum", v.(int)+4)
			ctx.Next()
		})

		g0.On("/t1", func(ctx *xray.Context) {
			v, _ := ctx.Get("sum")
			ctx.Set("sum", v.(int)+2)
			ctx.Next()
		}, func(ctx *xray.Context) {
			v, _ := ctx.Get("sum")
			as.Equal(3, v.(int))
			as.Equal(ctx.Writer.Protocol(), xray.ProtocolHTTP)

			{
				ctx.Writer.Header().Set(xray.ContentType, "plain/text")
				as.NoError(ctx.WriteString(http.StatusOK, "OK"))
				_, ok := ctx.Writer.Raw().(http.ResponseWriter)
				as.Equal(true, ok)
				as.Equal("plain/text", ctx.Writer.Header().Get(xray.ContentType))
			}
		})

		g0.On("t2", func(ctx *xray.Context) {
			v, _ := ctx.Get("sum")
			as.Equal(1, v.(int))
		})

		g1.On("t3", func(ctx *xray.Context) {
			v, _ := ctx.Get("sum")
			as.Equal(5, v.(int))
		})

		g2 := g0.Group("session")

		g2.On("t4", func(ctx *xray.Context) {
			v, _ := ctx.Get("sum")
			as.Equal(1, v.(int))
		})

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

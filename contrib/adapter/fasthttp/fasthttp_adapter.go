package fasthttp

import (
	"bytes"
	"github.com/lxzan/xray"
	"github.com/valyala/fasthttp"
	"io"
	"reflect"
	"unsafe"
)

type readCloser struct {
	payload []byte
	reader  io.Reader
}

func (c *readCloser) Close() error { return nil }

func (c *readCloser) Read(p []byte) (n int, err error) { return c.reader.Read(p) }

func (c *readCloser) Bytes() []byte { return c.payload }

type responseWriter struct {
	writer *fasthttp.Response
	header xray.Header
}

func (c *responseWriter) Write(p []byte) (n int, err error) {
	return c.writer.BodyWriter().Write(p)
}

func (c *responseWriter) Protocol() string {
	return xray.ProtocolHTTP
}

func (c *responseWriter) Header() xray.Header {
	return c.header
}

func (c *responseWriter) Code(code int) {
	c.writer.SetStatusCode(code)
}

func (c *responseWriter) Raw() any {
	return c.writer
}

func (c *responseWriter) Flush() error {
	return nil
}

func NewAdapter(router *xray.Router) *Adapter {
	return &Adapter{router: router}
}

// Adapter FastHTTP适配器
type Adapter struct {
	router *xray.Router
}

// SetRouter 设置路由器
func (c *Adapter) SetRouter(r *xray.Router) *Adapter {
	c.router = r
	return c
}

// ServeFastHTTP 服务HTTP
func (c *Adapter) ServeFastHTTP(ctx *fasthttp.RequestCtx) {
	var body = ctx.Request.Body()
	var r = &xray.Request{
		Raw:    &ctx.Request,
		Header: &requestHeader{RequestHeader: &ctx.Request.Header},
		Method: b2s(ctx.Method()),
		Body:   &readCloser{payload: body, reader: bytes.NewReader(body)},
	}

	var uctx = xray.NewContext(c.router, r, &responseWriter{
		writer: &ctx.Response,
		header: &responseHeader{ResponseHeader: &ctx.Response.Header},
	})
	c.router.EmitEvent(r.Method, b2s(ctx.Request.URI().Path()), uctx)
}

// b2s converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func b2s(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

// s2b converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func s2b(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

type requestHeader struct {
	*fasthttp.RequestHeader
}

func (c *requestHeader) New() xray.Header {
	return &requestHeader{}
}

func (c *requestHeader) Get(key string) string { return b2s(c.Peek(key)) }

func (c *requestHeader) Range(f func(key string, value string) bool) {
	next := true
	c.VisitAll(func(key, value []byte) {
		if !next {
			return
		}
		next = f(b2s(key), b2s(value))
	})
}

type responseHeader struct {
	*fasthttp.ResponseHeader
}

func (c *responseHeader) New() xray.Header { return &responseHeader{} }

func (c *responseHeader) Get(key string) string { return b2s(c.Peek(key)) }

func (c *responseHeader) Range(f func(key string, value string) bool) {
	next := true
	c.VisitAll(func(key, value []byte) {
		if !next {
			return
		}
		next = f(b2s(key), b2s(value))
	})
}

package fasthttp

import (
	"bytes"
	"github.com/lxzan/uRouter"
	"github.com/valyala/fasthttp"
	"reflect"
	"unsafe"
)

type responseWriter struct {
	writer *fasthttp.Response
	header uRouter.Header
}

func (c *responseWriter) Write(p []byte) (n int, err error) {
	return c.writer.BodyWriter().Write(p)
}

func (c *responseWriter) Protocol() string {
	return uRouter.ProtocolHTTP
}

func (c *responseWriter) Header() uRouter.Header {
	return c.header
}

func (c *responseWriter) Code(code int) {
	c.writer.SetStatusCode(code)
}

func (c *responseWriter) Raw() interface{} {
	return c.writer
}

func (c *responseWriter) Flush() error {
	return nil
}

func NewAdapter(router *uRouter.Router) *Adapter {
	return &Adapter{router: router}
}

// Adapter FastHTTP适配器
type Adapter struct {
	router *uRouter.Router
}

// SetRouter 设置路由器
func (c *Adapter) SetRouter(r *uRouter.Router) *Adapter {
	c.router = r
	return c
}

// ServeFastHTTP 服务HTTP
func (c *Adapter) ServeFastHTTP(fctx *fasthttp.RequestCtx) {
	var r = &uRouter.Request{
		Raw:    &fctx.Request,
		Header: &requestHeader{RequestHeader: &fctx.Request.Header},
		Action: b2s(fctx.Method()),
		Body:   bytes.NewBuffer(fctx.Request.Body()),
	}

	var ctx = uRouter.NewContext(r, &responseWriter{
		writer: &fctx.Response,
		header: &responseHeader{ResponseHeader: &fctx.Response.Header},
	})
	c.router.EmitEvent(r.Action, b2s(fctx.Request.URI().Path()), ctx)
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

func (c *requestHeader) Number() uint8 { return 0 }

func (c *requestHeader) Generate() uRouter.Header { return nil }

func (c *requestHeader) Close() {}

func (c *requestHeader) Get(key string) string {
	return b2s(c.Peek(key))
}

func (c *requestHeader) Range(f func(key string, value string)) {
	c.VisitAll(func(key, value []byte) {
		f(b2s(key), b2s(value))
	})
}

type responseHeader struct {
	*fasthttp.ResponseHeader
}

func (c *responseHeader) Number() uint8 { return 0 }

func (c *responseHeader) Generate() uRouter.Header { return nil }

func (c *responseHeader) Close() {}

func (c *responseHeader) Get(key string) string { return b2s(c.Peek(key)) }

func (c *responseHeader) Range(f func(key string, value string)) {
	c.VisitAll(func(key, value []byte) {
		f(b2s(key), b2s(value))
	})
}

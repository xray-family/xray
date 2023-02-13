package fasthttp

import (
	"bytes"
	"github.com/lxzan/uRouter"
	"github.com/valyala/fasthttp"
	"reflect"
	"sync"
	"unsafe"
)

var headerTemplate = new(uRouter.MapHeader)

type responseWriter struct {
	once   *sync.Once
	code   int
	writer *fasthttp.Response
	header uRouter.Header
}

func (c *responseWriter) Close() {
	c.header.Close()
	c.header = nil
}

func (c *responseWriter) Write(p []byte) (n int, err error) {
	var h = &c.writer.Header
	c.once.Do(func() {
		c.writer.SetStatusCode(c.code)
		c.header.Range(func(key, value string) {
			h.Set(key, value)
		})
	})
	return c.writer.BodyWriter().Write(p)
}

func (c *responseWriter) Protocol() string {
	return uRouter.ProtocolHTTP
}

func (c *responseWriter) Header() uRouter.Header {
	return c.header
}

func (c *responseWriter) Code(code int) {
	c.code = code
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

// Adapter HTTP适配器
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
		Header: headerTemplate.Generate(),
		Action: b2s(fctx.Method()),
		Body:   bytes.NewBuffer(fctx.Request.Body()),
	}
	fctx.Request.Header.VisitAll(func(key, value []byte) {
		r.Header.Set(b2s(key), b2s(value))
	})

	var ctx = uRouter.NewContext(r, &responseWriter{
		once:   &sync.Once{},
		writer: &fctx.Response,
		header: headerTemplate.Generate(),
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

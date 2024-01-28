package http

import (
	"github.com/lxzan/xray"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
}

func (c *responseWriter) Protocol() string {
	return xray.ProtocolHTTP
}

func (c *responseWriter) Header() xray.Header {
	return xray.NewHttpHeader(c.ResponseWriter.Header())
}

func (c *responseWriter) Code(code int) {
	c.ResponseWriter.WriteHeader(code)
}

func (c *responseWriter) Raw() any {
	return c.ResponseWriter
}

func (c *responseWriter) Flush() error {
	return nil
}

func NewAdapter(router *xray.Router) *Adapter {
	return &Adapter{router: router}
}

// Adapter HTTP适配器
type Adapter struct {
	router *xray.Router
}

// SetRouter 设置路由器
func (c *Adapter) SetRouter(r *xray.Router) *Adapter {
	c.router = r
	return c
}

// ServeHTTP 服务HTTP
func (c *Adapter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var r = &xray.Request{
		Raw:    request,
		Header: xray.NewHttpHeader(request.Header),
		Body:   request.Body,
		Action: request.Method,
	}
	var ctx = xray.NewContext(r, &responseWriter{ResponseWriter: writer})
	c.router.EmitEvent(r.Action, request.URL.Path, ctx)
}

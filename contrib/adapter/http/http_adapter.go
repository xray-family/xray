package http

import (
	"github.com/lxzan/uRouter"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
}

func (c *responseWriter) Protocol() string {
	return uRouter.ProtocolHTTP
}

func (c *responseWriter) Header() uRouter.Header {
	return uRouter.NewHttpHeader(c.ResponseWriter.Header())
}

func (c *responseWriter) Code(code int) {
	c.ResponseWriter.WriteHeader(code)
}

func (c *responseWriter) Raw() interface{} {
	return c.ResponseWriter
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

// ServeHTTP 服务HTTP
func (c *Adapter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var r = &uRouter.Request{
		Raw:    request,
		Header: uRouter.NewHttpHeader(request.Header),
		Body:   request.Body,
		Action: request.Method,
	}
	var ctx = uRouter.NewContext(r, &responseWriter{ResponseWriter: writer})
	c.router.EmitEvent(r.Action, request.URL.Path, ctx)
}

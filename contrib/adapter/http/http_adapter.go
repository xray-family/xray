package http

import (
	"github.com/xray-family/xray"
	"net/http"
)

type responseWriter struct{ http.ResponseWriter }

func (c *responseWriter) Protocol() string {
	return xray.ProtocolHTTP
}

func (c *responseWriter) Header() xray.Header {
	return &xray.HttpHeader{Header: c.ResponseWriter.Header()}
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

// ServeHTTP 服务HTTP
func (c *Adapter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var r = &xray.Request{
		Raw:    request,
		Header: &xray.HttpHeader{Header: request.Header},
		Body:   request.Body,
		Method: request.Method,
	}
	var ctx = xray.NewContext(c.router, r, &responseWriter{ResponseWriter: writer})
	c.router.EmitEvent(r.Method, request.URL.Path, ctx)
}

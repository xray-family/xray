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

func NewAdapter() *Adapter {
	return &Adapter{Router: uRouter.New()}
}

// Adapter HTTP适配器
type Adapter struct {
	*uRouter.Router
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
	c.Router.EmitAction(r.Action, request.URL.Path, ctx)
}

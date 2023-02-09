package http

import (
	"github.com/lxzan/uRouter"
	"github.com/lxzan/uRouter/codec"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
}

func (c *responseWriter) Protocol() string {
	return codec.ProtocolHTTP
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

func NewAdapter(r *uRouter.Router) *Adapter {
	return &Adapter{router: r}
}

// Adapter HTTP适配器
type Adapter struct {
	router *uRouter.Router
}

// ServeHTTP 服务HTTP
func (c *Adapter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var r = &uRouter.Request{
		Raw:    request,
		Header: uRouter.NewHttpHeader(request.Header),
		Body:   request.Body,
	}
	request.Header.Set(uRouter.XPath, request.URL.Path)
	var ctx = uRouter.NewContext(r, &responseWriter{ResponseWriter: writer})
	c.router.Emit(ctx)
}

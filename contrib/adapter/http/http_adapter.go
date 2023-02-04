package http

import (
	"github.com/lxzan/uRouter"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
}

func (c *responseWriter) Header() uRouter.Header {
	return uRouter.NewHttpHeader(c.ResponseWriter.Header())
}

func (c *responseWriter) Code(code int) {
	c.ResponseWriter.WriteHeader(code)
}

func (c *responseWriter) RawResponseWriter() interface{} {
	return c.ResponseWriter
}

func (c *responseWriter) Flush() error {
	return nil
}

func NewAdapter(r *uRouter.Router) *Adapter {
	return &Adapter{router: r}
}

type Adapter struct {
	router *uRouter.Router
}

func (c *Adapter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var r = &uRouter.Request{
		RawRequest: request,
		Header:     uRouter.NewHttpHeader(request.Header),
		Body:       request.Body,
	}
	request.Header.Set(uRouter.XPath, request.URL.Path)
	var ctx = uRouter.NewContext(r, &responseWriter{ResponseWriter: writer})
	c.router.Emit(ctx)
}

package http

import (
	"github.com/lxzan/uRouter"
	"net/http"
)

type Writer struct {
	http.ResponseWriter
}

func (c *Writer) Header() uRouter.Header {
	return uRouter.NewHttpHeader(c.ResponseWriter.Header())
}

func (c *Writer) Code(code int) {
	c.ResponseWriter.WriteHeader(code)
}

func (c *Writer) RawResponseWriter() interface{} {
	return c.ResponseWriter
}

func (c *Writer) Flush() error {
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
	var ctx = uRouter.NewContext(r, &Writer{ResponseWriter: writer})
	c.router.Emit(ctx)
}

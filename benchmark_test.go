package xray

import (
	"net/http"
	"testing"
)

var anyMethods = []string{
	http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
	http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
	http.MethodTrace,
}

func ServeAny(r *Router, path string, hs ...HandlerFunc) {
	for _, method := range anyMethods {
		r.OnEvent(method, path, hs...)
	}
}

func BenchmarkOneRoute(B *testing.B) {
	router := New(WithGreeting(false, 0))
	router.GET("/ping", func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkRecoveryMiddleware(B *testing.B) {
	router := New(WithGreeting(false, 0))
	router.Use(Recovery())
	router.GET("/", func(c *Context) {})
	runRequest(B, router, "GET", "/")
}

func BenchmarkManyHandlers(B *testing.B) {
	router := New(WithGreeting(false, 0))
	router.Use(Recovery())
	router.Use(func(c *Context) {})
	router.Use(func(c *Context) {})
	router.GET("/ping", func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func Benchmark5Params(B *testing.B) {
	router := New(WithGreeting(false, 0))
	router.Use(func(c *Context) {})
	router.GET("/param/:param1/:params2/:param3/:param4/:param5", func(c *Context) {})
	runRequest(B, router, "GET", "/param/path/to/parameter/john/12345")
}

func BenchmarkOneRouteJSON(B *testing.B) {
	router := New(WithGreeting(false, 0))
	data := struct {
		Status string `json:"status"`
	}{"ok"}
	router.GET("/json", func(c *Context) {
		c.WriteJSON(http.StatusOK, data)
	})
	runRequest(B, router, "GET", "/json")
}

func BenchmarkOneRouteSet(B *testing.B) {
	router := New(WithGreeting(false, 0))
	router.GET("/ping", func(c *Context) {
		c.Set("key", "value")
	})
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkOneRouteString(B *testing.B) {
	router := New(WithGreeting(false, 0))
	router.GET("/text", func(c *Context) {
		c.WriteString(http.StatusOK, "this is a plain text")
	})
	runRequest(B, router, "GET", "/text")
}

func BenchmarkManyRoutesFist(B *testing.B) {
	router := New(WithGreeting(false, 0))
	ServeAny(router, "/ping", func(c *Context) {})
	runRequest(B, router, "GET", "/ping")
}

func BenchmarkManyRoutesLast(B *testing.B) {
	router := New(WithGreeting(false, 0))
	ServeAny(router, "/ping", func(c *Context) {})
	runRequest(B, router, "OPTIONS", "/ping")
}

func Benchmark404(B *testing.B) {
	router := New(WithGreeting(false, 0))
	ServeAny(router, "/something", func(c *Context) {})
	router.SetHandlerNotFound(func(ctx *Context) {

	})
	runRequest(B, router, "GET", "/ping")
}

func Benchmark404Many(B *testing.B) {
	router := New(WithGreeting(false, 0))
	router.GET("/", func(c *Context) {})
	router.GET("/path/to/something", func(c *Context) {})
	router.GET("/post/:id", func(c *Context) {})
	router.GET("/view/:id", func(c *Context) {})
	router.GET("/favicon.ico", func(c *Context) {})
	router.GET("/robots.txt", func(c *Context) {})
	router.GET("/delete/:id", func(c *Context) {})
	router.GET("/user/:id/:mode", func(c *Context) {})

	router.SetHandlerNotFound(func(ctx *Context) {

	})
	runRequest(B, router, "GET", "/viewfake")
}

type mockWriter struct {
	headers http.Header
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		http.Header{},
	}
}

func (m *mockWriter) Header() (h http.Header) {
	return m.headers
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockWriter) WriteHeader(int) {}

func runRequest(B *testing.B, router *Router, method, path string) {
	req, _ := http.NewRequest(method, path, nil)
	var r = &Request{
		Raw:    req,
		Header: &HttpHeader{Header: req.Header},
		Body:   req.Body,
		Method: req.Method,
	}
	var ctx = NewContext(router, r, newResponseWriterMocker())
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		router.EmitEvent(method, path, ctx)
	}
}

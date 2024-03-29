package xray

import (
	"net/http"
	"testing"
)

func BenchmarkOneRoute(b *testing.B) {
	var option = WithGreeting(false, 0)
	router := New(option)
	router.GET("/ping", func(c *Context) {
	})
	ctx := newContextMocker(option)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/ping", ctx)
	}
}

func BenchmarkOneRouteDynamic(b *testing.B) {
	option := WithGreeting(false, 0)
	router := New(option)
	router.GET("/user/:id", func(c *Context) {
	})
	ctx := newContextMocker(option)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/user/1", ctx)
	}
}

func BenchmarkRecoveryMiddleware(b *testing.B) {
	option := WithGreeting(false, 0)
	router := New(option)
	router.Use(Recovery())
	router.GET("/", func(c *Context) {
	})
	ctx := newContextMocker(option)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/", ctx)
	}
}

func Benchmark5Params(b *testing.B) {
	option := WithGreeting(false, 0)
	router := New(option)
	router.Use(func(ctx *Context) {})
	router.GET("/param/:param1/:params2/:param3/:param4/:param5", func(c *Context) {
	})
	ctx := newContextMocker(option)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/param/path/to/parameter/john/12345", ctx)
	}
}

func BenchmarkOneRouteJSON(b *testing.B) {
	option := WithGreeting(false, 0)
	router := New(option)
	router.Use(func(ctx *Context) {})
	data := struct {
		Status string `json:"status"`
	}{"ok"}
	router.GET("/json", func(c *Context) {
		defer c.Request.Body.Close()
		_ = c.WriteJSON(http.StatusOK, data)
	})
	ctx := newContextMocker(option)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/json", ctx)
	}
}

func Benchmark404(b *testing.B) {
	option := WithGreeting(false, 0)
	router := New(option)
	router.GET("/", func(c *Context) {})
	router.GET("/something", func(c *Context) {})
	ctx := newContextMocker(option)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/ping", ctx)
	}
}

func Benchmark404Many(b *testing.B) {
	option := WithGreeting(false, 0)
	router := New(option)
	router.GET("/", func(c *Context) {})
	router.GET("/path/to/something", func(c *Context) {})
	router.GET("/post/:id", func(c *Context) {})
	router.GET("/view/:id", func(c *Context) {})
	router.GET("/favicon.ico", func(c *Context) {})
	router.GET("/robots.txt", func(c *Context) {})
	router.GET("/delete/:id", func(c *Context) {})
	router.GET("/user/:id/:mode", func(c *Context) {})
	ctx := newContextMocker(option)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/viewfake", ctx)
	}
}

package fasthttp

import (
	"github.com/valyala/fasthttp"
	"github.com/xray-family/xray"
	"testing"
)

func BenchmarkAdapter_ServeFastHTTP(b *testing.B) {
	var ctx = &fasthttp.RequestCtx{
		Request:  fasthttp.Request{},
		Response: fasthttp.Response{},
	}
	var router = xray.New()
	router.GET("/", func(ctx *xray.Context) {})
	adapter := NewAdapter(router)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		adapter.ServeFastHTTP(ctx)
	}
}

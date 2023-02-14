package fasthttp

import (
	"github.com/lxzan/uRouter"
	"github.com/valyala/fasthttp"
	"testing"
)

func BenchmarkAdapter_ServeFastHTTP(b *testing.B) {
	var ctx = &fasthttp.RequestCtx{
		Request:  fasthttp.Request{},
		Response: fasthttp.Response{},
	}
	var router = uRouter.New()
	router.OnGET("/", func(ctx *uRouter.Context) {})
	router.StartSilently()
	adapter := NewAdapter(router)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		adapter.ServeFastHTTP(ctx)
	}
}

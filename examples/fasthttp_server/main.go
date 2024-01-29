package main

import (
	"github.com/lxzan/xray"
	fasthttpAdapter "github.com/lxzan/xray/contrib/adapter/fasthttp"
	"github.com/valyala/fasthttp"
	"net/http"
)

func main() {
	var router = xray.New()
	router.Use(xray.Recovery())

	router.OnGET("/test", func(ctx *xray.Context) {
		_ = ctx.WriteBytes(http.StatusOK, []byte("hello"))
	})

	// 开启一个原生的fasthttp server, 对比下性能
	go func() {
		if err := fasthttp.ListenAndServe(":3001", func(ctx *fasthttp.RequestCtx) {
			ctx.SetBody([]byte("hello"))
		}); err != nil {
			router.Logger().Panic(err.Error())
		}
	}()

	// 开启一个原生的http server, 对比下性能
	go func() {
		if err := http.ListenAndServe(":3002", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte("hello"))
		})); err != nil {
			router.Logger().Panic(err.Error())
		}
	}()

	if err := fasthttp.ListenAndServe(":3000", fasthttpAdapter.NewAdapter(router).ServeFastHTTP); err != nil {
		router.Logger().Panic(err.Error())
	}
}

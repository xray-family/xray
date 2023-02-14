package main

import (
	"github.com/lxzan/uRouter"
	fasthttpAdapter "github.com/lxzan/uRouter/contrib/adapter/fasthttp"
	"github.com/valyala/fasthttp"
	"net/http"
)

func main() {
	var router = uRouter.New()
	router.Use(uRouter.Recovery())

	router.OnGET("/test", func(ctx *uRouter.Context) {
		_ = ctx.WriteBytes(http.StatusOK, []byte("hello"))
	})

	router.Start()

	// 开启一个原生的fasthttp server, 对比下性能
	go func() {
		if err := fasthttp.ListenAndServe(":3001", func(ctx *fasthttp.RequestCtx) {
			ctx.SetBody([]byte("hello"))
		}); err != nil {
			uRouter.Logger().Panic(err.Error())
		}
	}()

	// 开启一个原生的http server, 对比下性能
	go func() {
		if err := http.ListenAndServe(":3002", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte("hello"))
		})); err != nil {
			uRouter.Logger().Panic(err.Error())
		}
	}()

	if err := fasthttp.ListenAndServe(":3000", fasthttpAdapter.NewAdapter(router).ServeFastHTTP); err != nil {
		uRouter.Logger().Panic(err.Error())
	}
}

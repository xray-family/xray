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

	router.OnGET("/test", Test)

	router.Start()

	// 开启一个原生的server, 对比下性能损失
	go func() {
		if err := fasthttp.ListenAndServe(":3001", func(ctx *fasthttp.RequestCtx) {
			ctx.SetBody([]byte("hello"))
		}); err != nil {
			uRouter.Logger().Panic(err.Error())
		}
	}()

	if err := fasthttp.ListenAndServe(":3000", fasthttpAdapter.NewAdapter(router).ServeFastHTTP); err != nil {
		uRouter.Logger().Panic(err.Error())
	}
}

func Test(ctx *uRouter.Context) {
	_ = ctx.WriteString(http.StatusOK, "hello")
}

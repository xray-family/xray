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

	if err := fasthttp.ListenAndServe(":3000", fasthttpAdapter.NewAdapter(router).ServeFastHTTP); err != nil {
		uRouter.Logger().Panic(err.Error())
	}
}

func Test(ctx *uRouter.Context) {
	_ = ctx.WriteString(http.StatusOK, "hello")
	ctx.Close()
}

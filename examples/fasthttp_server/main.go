package main

import (
	"github.com/valyala/fasthttp"
	"github.com/xray-family/xray"
	fasthttpAdapter "github.com/xray-family/xray/contrib/adapter/fasthttp"
	"net/http"
)

func main() {
	var router = xray.New()
	router.Use(xray.Recovery())

	router.GET("/api/v1/test", func(ctx *xray.Context) {
		_ = ctx.WriteBytes(http.StatusOK, []byte("hello"))
	})

	if err := fasthttp.ListenAndServe(":3002", fasthttpAdapter.NewAdapter(router).ServeFastHTTP); err != nil {
		router.Logger().Panic(err.Error())
	}
}

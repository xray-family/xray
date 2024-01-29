package main

import (
	"github.com/lxzan/xray"
	httpAdapter "github.com/lxzan/xray/contrib/adapter/http"
	"github.com/lxzan/xray/contrib/log/zerolog"
	"github.com/quic-go/quic-go/http3"
	"net/http"
)

func main() {
	var router = xray.New(xray.WithLogger(zerolog.Logger))
	router.Use(xray.Recovery())
	var group = router.Group("/api/v1")

	group.OnEvent(http.MethodGet, "/user/:uid/article/:aid", func(ctx *xray.Context) {
		_ = ctx.WriteJSON(http.StatusOK, xray.Any{
			"uid": ctx.Param("uid"),
			"aid": ctx.Param("aid"),
		})
	})

	if err := http3.ListenAndServe(":3000", "cert/server.pem", "cert/server-key.pem", httpAdapter.NewAdapter(router)); err != nil {
		router.Logger().Panic(err.Error())
	}
}

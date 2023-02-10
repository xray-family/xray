package main

import (
	"github.com/lxzan/uRouter"
	httpAdapter "github.com/lxzan/uRouter/contrib/adapter/http"
	"github.com/lxzan/uRouter/contrib/codec/jsoniter"
	"github.com/lxzan/uRouter/contrib/log/zerolog"
	"github.com/quic-go/quic-go/http3"
	"net/http"
)

func init() {
	uRouter.SetJsonCodec(jsoniter.JsoniterCodec)
	uRouter.SetLogger(zerolog.ZeroLogger)
}

func main() {
	var router = uRouter.New()
	router.Use(uRouter.Recovery())
	var group = router.Group("/api/v1")

	group.OnEvent(http.MethodGet, "/user/:uid/article/:aid", func(ctx *uRouter.Context) {
		_ = ctx.WriteJSON(http.StatusOK, uRouter.Any{
			"uid": ctx.Param("uid"),
			"aid": ctx.Param("aid"),
		})
	})

	router.Display()
	if err := http3.ListenAndServe(":3000", "cert/server.pem", "cert/server-key.pem", httpAdapter.NewAdapter(router)); err != nil {
		uRouter.Logger().Panic(err.Error())
	}
}

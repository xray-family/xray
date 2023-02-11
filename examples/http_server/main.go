package main

import (
	"github.com/lxzan/uRouter"
	httpAdapter "github.com/lxzan/uRouter/contrib/adapter/http"
	"github.com/lxzan/uRouter/contrib/codec/jsoniter"
	"github.com/lxzan/uRouter/contrib/log/zerolog"
	"net/http"
)

func init() {
	uRouter.SetLogger(zerolog.ZeroLogger)
	uRouter.SetJsonCodec(jsoniter.JsoniterCodec)
}

func main() {
	r := uRouter.New()
	r.Use(uRouter.Recovery(), uRouter.AccessLog())
	group := r.Group("/api/v1")
	NewController().Mapping(group)
	r.Start()

	if err := http.ListenAndServe(":3000", httpAdapter.NewAdapter(r)); err != nil {
		uRouter.Logger().Panic(err.Error())
	}
}

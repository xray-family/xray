package main

import (
	"github.com/lxzan/xray"
	httpAdapter "github.com/lxzan/xray/contrib/adapter/http"
	"github.com/lxzan/xray/contrib/codec/jsoniter"
	"github.com/lxzan/xray/contrib/log/zerolog"
	"net/http"
)

func init() {
	xray.SetLogger(zerolog.ZeroLogger)
	xray.SetJsonCodec(jsoniter.JsoniterCodec)
}

func main() {
	r := xray.New()
	r.Use(xray.Recovery(), xray.AccessLog())
	group := r.Group("/api/v1")
	NewController().Mapping(group)

	if err := http.ListenAndServe(":3000", httpAdapter.NewAdapter(r)); err != nil {
		xray.Logger().Panic(err.Error())
	}
}

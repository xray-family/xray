package main

import (
	"github.com/lxzan/xray"
	httpAdapter "github.com/lxzan/xray/contrib/adapter/http"
	"github.com/lxzan/xray/contrib/log/zerolog"
	"net/http"
)

func main() {
	r := xray.New(xray.WithLogger(zerolog.Logger))
	r.Use(xray.Recovery(), xray.AccessLog())
	group := r.Group("/api/v1")
	NewController().Mapping(group)

	if err := http.ListenAndServe(":3000", httpAdapter.NewAdapter(r)); err != nil {
		r.Logger().Panic(err.Error())
	}
}

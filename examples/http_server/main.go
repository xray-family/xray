package main

import (
	"github.com/xray-family/xray"
	httpAdapter "github.com/xray-family/xray/contrib/adapter/http"
	"github.com/xray-family/xray/contrib/log/zerolog"
	"net/http"
)

func main() {
	r := xray.New(xray.WithLogger(zerolog.Logger))
	r.Use(xray.Recovery())
	group := r.Group("/api/v1")
	NewController().Mapping(group)

	if err := http.ListenAndServe(":3001", httpAdapter.NewAdapter(r)); err != nil {
		r.Logger().Panic(err.Error())
	}
}

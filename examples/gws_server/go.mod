module github.com/lxzan/xray/examples/gws_server

go 1.18

replace (
	github.com/lxzan/xray => ../../
	github.com/lxzan/xray/contrib/adapter/gws => ../../contrib/adapter/gws
	github.com/lxzan/xray/contrib/adapter/http => ../../contrib/adapter/http
	github.com/lxzan/xray/contrib/log/zerolog => ../../contrib/log/zerolog
)

require (
	github.com/lxzan/gws v1.8.0
	github.com/lxzan/xray v1.0.0
	github.com/lxzan/xray/contrib/adapter/gws v0.0.0-00010101000000-000000000000
	github.com/lxzan/xray/contrib/adapter/http v1.0.0
	github.com/lxzan/xray/contrib/codec/jsoniter v0.0.0-20230208023827-6baa50a514b5
	github.com/lxzan/xray/contrib/log/zerolog v0.0.0-00010101000000-000000000000
)

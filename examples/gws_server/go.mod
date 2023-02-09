module github.com/lxzan/uRouter/examples/gws_server

go 1.16

replace (
	github.com/lxzan/uRouter => ../../
	github.com/lxzan/uRouter/contrib/adapter/gws => ../../contrib/adapter/gws
	github.com/lxzan/uRouter/contrib/adapter/http => ../../contrib/adapter/http
	github.com/lxzan/uRouter/contrib/log/zerolog => ../../contrib/log/zerolog
)

require (
	github.com/lxzan/gws v1.2.14
	github.com/lxzan/uRouter v1.0.0
	github.com/lxzan/uRouter/contrib/adapter/gws v0.0.0-00010101000000-000000000000
	github.com/lxzan/uRouter/contrib/adapter/http v1.0.0
	github.com/lxzan/uRouter/contrib/codec/jsoniter v0.0.0-20230208023827-6baa50a514b5
	github.com/lxzan/uRouter/contrib/log/zerolog v0.0.0-00010101000000-000000000000
)

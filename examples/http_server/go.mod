module github.com/lxzan/xray/examples/http_server

go 1.18

replace (
	github.com/lxzan/xray => ../../
	github.com/lxzan/xray/contrib/adapter/http => ../../contrib/adapter/http
	github.com/lxzan/xray/contrib/codec/jsoniter => ../../contrib/codec/jsoniter
	github.com/lxzan/xray/contrib/log/zerolog => ../../contrib/log/zerolog
)

require (
	github.com/lxzan/xray v1.0.0
	github.com/lxzan/xray/contrib/adapter/http v1.0.0
	github.com/lxzan/xray/contrib/codec/jsoniter v0.0.0-20230208023827-6baa50a514b5
	github.com/lxzan/xray/contrib/log/zerolog v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10 // indirect
)

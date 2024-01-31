module github.com/xray-family/xray/examples/gws_server

go 1.18

replace (
	github.com/xray-family/xray => ../../
	github.com/xray-family/xray/contrib/adapter/gws => ../../contrib/adapter/gws
	github.com/xray-family/xray/contrib/adapter/http => ../../contrib/adapter/http
	github.com/xray-family/xray/contrib/codec/jsoniter => ../../contrib/codec/jsoniter
	github.com/xray-family/xray/contrib/log/zerolog => ../../contrib/log/zerolog
)

require (
	github.com/lxzan/gws v1.8.0
	github.com/xray-family/xray v1.0.0
	github.com/xray-family/xray/contrib/adapter/gws v0.0.0-00010101000000-000000000000
	github.com/xray-family/xray/contrib/adapter/http v1.0.0
	github.com/xray-family/xray/contrib/codec/jsoniter v0.0.0-20240129145014-d89ba3b3f26a
	github.com/xray-family/xray/contrib/log/zerolog v0.0.0-00010101000000-000000000000
)

require (
	github.com/dolthub/maphash v0.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.5 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/rs/zerolog v1.29.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
)

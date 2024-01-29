module github.com/lxzan/xray/examples/fasthttp_server

go 1.18

replace (
	github.com/lxzan/xray => ../../
	github.com/lxzan/xray/contrib/adapter/fasthttp => ../../contrib/adapter/fasthttp
)

require (
	github.com/lxzan/xray v1.0.0
	github.com/lxzan/xray/contrib/adapter/fasthttp v0.0.0-00010101000000-000000000000
	github.com/valyala/fasthttp v1.44.0
)

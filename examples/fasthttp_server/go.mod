module github.com/lxzan/uRouter/examples/fasthttp_server

go 1.16

replace (
	github.com/lxzan/uRouter => ../../
	github.com/lxzan/uRouter/contrib/adapter/fasthttp => ../../contrib/adapter/fasthttp
)

require (
	github.com/lxzan/uRouter v1.0.0
	github.com/lxzan/uRouter/contrib/adapter/fasthttp v0.0.0-00010101000000-000000000000
	github.com/valyala/fasthttp v1.44.0
)

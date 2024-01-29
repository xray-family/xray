module github.com/xray-family/xray/examples/fasthttp_server

go 1.18

replace (
	github.com/xray-family/xray => ../../
	github.com/xray-family/xray/contrib/adapter/fasthttp => ../../contrib/adapter/fasthttp
)

require (
	github.com/valyala/fasthttp v1.44.0
	github.com/xray-family/xray v0.0.1
	github.com/xray-family/xray/contrib/adapter/fasthttp v0.0.0-00010101000000-000000000000
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/klauspost/compress v1.17.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
)

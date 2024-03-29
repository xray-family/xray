module github.com/xray-family/xray/examples/http_server

go 1.18

replace (
	github.com/xray-family/xray => ../../
	github.com/xray-family/xray/contrib/adapter/http => ../../contrib/adapter/http
	github.com/xray-family/xray/contrib/log/zerolog => ../../contrib/log/zerolog
)

require (
	github.com/xray-family/xray v0.0.1
	github.com/xray-family/xray/contrib/adapter/http v1.0.0
	github.com/xray-family/xray/contrib/log/zerolog v0.0.0-00010101000000-000000000000
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/rs/zerolog v1.29.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
)

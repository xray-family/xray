module github.com/lxzan/xray/contrib/doc/swagger

go 1.18

replace (
	github.com/lxzan/xray => ../../../
	github.com/lxzan/xray/contrib/adapter/http => ../../../contrib/adapter/http
)

require (
	github.com/lxzan/xray v1.0.0
	github.com/lxzan/xray/contrib/adapter/http v0.0.0-20230214121206-a6001cd774ce
	github.com/stretchr/testify v1.8.1
	github.com/swaggo/files v1.0.0
	github.com/swaggo/swag v1.8.10
	golang.org/x/net v0.7.0
)

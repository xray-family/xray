module github.com/lxzan/uRouter/contrib/doc/swagger

go 1.16

replace (
	github.com/lxzan/uRouter => ../../../
	github.com/lxzan/uRouter/contrib/adapter/http => ../../../contrib/adapter/http
)

require (
	github.com/lxzan/uRouter v1.0.0
	github.com/lxzan/uRouter/contrib/adapter/http v0.0.0-20230214121206-a6001cd774ce
	github.com/stretchr/testify v1.8.1
	github.com/swaggo/files v1.0.0
	github.com/swaggo/swag v1.8.10
	golang.org/x/net v0.7.0
)

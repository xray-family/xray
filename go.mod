module github.com/lxzan/uRouter

go 1.16

replace (
	github.com/lxzan/uRouter/contrib/adapter/gws => ./contrib/adapter/gws
	github.com/lxzan/uRouter/contrib/adapter/http => ./contrib/adapter/http
)

require (
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/lxzan/gws v1.2.13
	github.com/lxzan/uRouter/contrib/adapter/gws v0.0.0-20230204030007-1bda8d97a560
	github.com/lxzan/uRouter/contrib/adapter/http v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.1
	github.com/ugorji/go/codec v1.2.9
)

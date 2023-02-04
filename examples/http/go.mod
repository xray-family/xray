module demo

go 1.16

replace (
github.com/lxzan/uRouter/contrib/adapter/http => ../../contrib/adapter/http
github.com/lxzan/uRouter  => ../../
)

require (

	github.com/lxzan/uRouter v0.0.0
	github.com/lxzan/uRouter/contrib/adapter/http v0.0.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

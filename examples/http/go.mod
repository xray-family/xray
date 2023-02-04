module github.com/lxzan/uRouter/examples/http

go 1.16

replace github.com/lxzan/uRouter/contrib/adapter/http => ../../contrib/adapter/http

require (
	github.com/lxzan/uRouter v1.0.0
	github.com/lxzan/uRouter/contrib/adapter/http v1.0.0
	github.com/ugorji/go/codec v1.2.9 // indirect
)

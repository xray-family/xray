module github.com/lxzan/uRouter/contrib/adapter/gws

go 1.16

replace (
	github.com/lxzan/gws v1.2.13 => ../../../../gws
	github.com/lxzan/uRouter => ../../../
)

require (
	github.com/lxzan/gws v1.2.13
	github.com/lxzan/uRouter v1.0.0
	github.com/stretchr/testify v1.8.1
)

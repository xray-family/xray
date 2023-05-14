package helper

import (
	"github.com/lxzan/uRouter"
	_ "unsafe"
)

//go:linkname DisplayRouter github.com/lxzan/uRouter.displayRouter
func DisplayRouter(c *uRouter.Router)

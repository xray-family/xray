package gws

import (
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
	"testing"
)

func TestRootController_GetMessage(t *testing.T) {
	var ctl = new(RootController)
	var ctx = uRouter.NewContext(
		&uRouter.Request{Raw: &gws.Message{}},
		nil,
	)
	ctl.GetMessage(ctx)
}

func TestRootController_GetConnection(t *testing.T) {
	var ctl = new(RootController)
	var ctx = uRouter.NewContext(
		&uRouter.Request{Raw: &gws.Message{}},
		&responseWriter{conn: &gws.Conn{}},
	)
	ctl.GetConnection(ctx)
}

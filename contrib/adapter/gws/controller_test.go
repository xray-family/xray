package gws

import (
	"github.com/lxzan/gws"
	"github.com/lxzan/xray"
	"testing"
)

func TestRootController_GetMessage(t *testing.T) {
	var ctl = new(RootController)
	var ctx = xray.NewContext(
		&xray.Request{Raw: &gws.Message{}},
		nil,
	)
	ctl.GetMessage(ctx)
}

func TestRootController_GetConnection(t *testing.T) {
	var ctl = new(RootController)
	var ctx = xray.NewContext(
		&xray.Request{Raw: &gws.Message{}},
		&responseWriter{conn: &gws.Conn{}},
	)
	ctl.GetConnection(ctx)
}

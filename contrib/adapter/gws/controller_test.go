package gws

import (
	"github.com/lxzan/gws"
	"github.com/xray-family/xray"
	"testing"
)

func TestRootController_GetMessage(t *testing.T) {
	var ctl = new(RootController)
	var router = xray.New()
	var ctx = xray.NewContext(
		router,
		&xray.Request{Raw: &gws.Message{}},
		nil,
	)
	ctl.GetMessage(ctx)
}

func TestRootController_GetConnection(t *testing.T) {
	var ctl = new(RootController)
	var router = xray.New()
	var ctx = xray.NewContext(
		router,
		&xray.Request{Raw: &gws.Message{}},
		&responseWriter{conn: &gws.Conn{}},
	)
	ctl.GetConnection(ctx)
}

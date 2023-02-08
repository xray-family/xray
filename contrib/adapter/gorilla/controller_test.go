package gorilla

import (
	"github.com/gorilla/websocket"
	"github.com/lxzan/uRouter"
	"testing"
)

func TestRootController_GetMessage(t *testing.T) {
	var ctl = new(RootController)
	var ctx = uRouter.NewContext(
		&uRouter.Request{Raw: &Message{}},
		nil,
	)
	ctl.GetMessage(ctx)
}

func TestRootController_GetConnection(t *testing.T) {
	var ctl = new(RootController)
	var ctx = uRouter.NewContext(
		&uRouter.Request{Raw: &Message{}},
		&responseWriter{conn: &websocket.Conn{}},
	)
	ctl.GetConnection(ctx)
}

package gorilla

import (
	"github.com/gorilla/websocket"
	"github.com/lxzan/uRouter"
)

type RootController struct{}

// GetMessage 获取websocket消息
func (c *RootController) GetMessage(ctx *uRouter.Context) *Message {
	return ctx.Request.Raw.(*Message)
}

// GetConnection 获取websocket连接
func (c *RootController) GetConnection(ctx *uRouter.Context) *websocket.Conn {
	return ctx.Writer.Raw().(*websocket.Conn)
}

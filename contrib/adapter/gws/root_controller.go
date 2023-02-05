package gws

import (
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
)

type RootController struct{}

// GetMessage 获取websocket消息
func (c *RootController) GetMessage(ctx *uRouter.Context) *gws.Message {
	return ctx.Request.Raw.(*gws.Message)
}

// GetConnection 获取websocket连接
func (c *RootController) GetConnection(ctx *uRouter.Context) *gws.Conn {
	return ctx.Writer.Raw().(*gws.Conn)
}

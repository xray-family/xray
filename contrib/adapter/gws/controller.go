package gws

import (
	"github.com/lxzan/gws"
	"github.com/xray-family/xray"
)

type RootController struct{}

// GetMessage 获取websocket消息
func (c *RootController) GetMessage(ctx *xray.Context) *gws.Message {
	return ctx.Request.Raw.(*gws.Message)
}

// GetConnection 获取websocket连接
func (c *RootController) GetConnection(ctx *xray.Context) *gws.Conn {
	return ctx.Writer.Raw().(*gws.Conn)
}

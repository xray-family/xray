package http

import (
	"github.com/lxzan/uRouter"
	"net/http"
)

type RootController struct{}

// GetRequest 获取http请求
func (c *RootController) GetRequest(ctx *uRouter.Context) *http.Request {
	return ctx.Request.Raw.(*http.Request)
}

// GetResponseWriter 获取http响应写入器
func (c *RootController) GetResponseWriter(ctx *uRouter.Context) http.ResponseWriter {
	return ctx.Writer.Raw().(http.ResponseWriter)
}

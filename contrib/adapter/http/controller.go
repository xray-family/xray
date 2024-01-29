package http

import (
	"github.com/lxzan/xray"
	"net/http"
	"strings"
)

const defaultIP = "127.0.0.1"

type RootController struct{}

// GetRequest 获取http请求
func (c *RootController) GetRequest(ctx *xray.Context) *http.Request {
	return ctx.Request.Raw.(*http.Request)
}

// GetResponseWriter 获取http响应写入器
func (c *RootController) GetResponseWriter(ctx *xray.Context) http.ResponseWriter {
	return ctx.Writer.Raw().(http.ResponseWriter)
}

// ClientIP 获取客户端IP, 优先从X-Real-IP中获取
// get the client IP, preferably from the X-Real-IP
func (c *RootController) ClientIP(ctx *xray.Context) string {
	if v := ctx.Request.Header.Get(xray.XRealIP); v != "" {
		return v
	}
	list := strings.Split(c.GetRequest(ctx).RemoteAddr, ":")
	if len(list) == 2 {
		return list[0]
	}
	return defaultIP
}

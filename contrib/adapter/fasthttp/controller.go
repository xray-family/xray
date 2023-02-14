package fasthttp

import (
	"github.com/lxzan/uRouter"
	"github.com/lxzan/uRouter/constant"
	"github.com/valyala/fasthttp"
	"strings"
)

const defaultIP = "127.0.0.1"

type RootController struct{}

// GetRequest 获取fasthttp请求
func (c *RootController) GetRequest(ctx *uRouter.Context) *fasthttp.Request {
	return ctx.Request.Raw.(*fasthttp.Request)
}

// GetResponseWriter 获取fasthttp响应写入器
func (c *RootController) GetResponseWriter(ctx *uRouter.Context) *fasthttp.Response {
	return ctx.Writer.Raw().(*fasthttp.Response)
}

// ClientIP 获取客户端IP, 优先从X-Real-IP中获取
// get the client IP, preferably from the X-Real-IP
func (c *RootController) ClientIP(ctx *uRouter.Context) string {
	if v := ctx.Request.Header.Get(constant.XRealIP); v != "" {
		return v
	}
	list := strings.Split(c.GetResponseWriter(ctx).RemoteAddr().String(), ":")
	if len(list) == 2 {
		return list[0]
	}
	return defaultIP
}

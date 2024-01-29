package xray

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// AccessLog 访问日志
func AccessLog() HandlerFunc {
	return func(ctx *Context) {
		var startTime = time.Now()
		ctx.Next()

		ctx.conf.logger.Debug(
			"access: protocol=%s, method=%s, path=%s, cost=%s",
			ctx.Writer.Protocol(),
			ctx.Request.Method,
			ctx.Request.RPath,
			time.Since(startTime).String(),
		)
	}
}

// Recovery 恢复模式, 防止服务崩溃
func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if fatalError := recover(); fatalError != nil {
				var msg = make([]byte, 0, 256)
				msg = append(msg, fmt.Sprintf("fatal error: %v\n", fatalError)...)
				for i := 1; true; i++ {
					_, caller, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					if !strings.Contains(caller, "src/runtime") {
						msg = append(msg, fmt.Sprintf("caller: %s, line: %d\n", caller, line)...)
					}
				}
				ctx.conf.logger.Info(string(msg))

				if ctx.Writer.Protocol() == ProtocolHTTP {
					_ = ctx.WriteString(http.StatusInternalServerError, "internal server error")
				}
			}
		}()
		ctx.Next()
	}
}

// HttpRequired 定义了HTTP协议允许通过的请求方法
func HttpRequired(methods ...string) HandlerFunc {
	for i, v := range methods {
		methods[i] = strings.ToUpper(v)
	}
	return func(ctx *Context) {
		if ctx.Writer.Protocol() != ProtocolHTTP {
			return
		}

		r, _ := ctx.Request.Raw.(*http.Request)
		for _, v := range methods {
			if r.Method == v {
				ctx.Next()
				return
			}
		}

		_ = ctx.WriteString(http.StatusForbidden, "method not allowed")
	}
}

// WebSocketRequired 只允许WebSocket协议请求通过
func WebSocketRequired() HandlerFunc {
	return func(ctx *Context) {
		if ctx.Writer.Protocol() == ProtocolWebSocket {
			ctx.Next()
		}
	}
}

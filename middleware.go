package uRouter

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// Protocol 协议过滤, 只允许一种协议访问
func Protocol(protocol string) HandlerFunc {
	return func(ctx *Context) {
		if ctx.Writer.Protocol() == protocol {
			ctx.Next()
			return
		}
	}
}

// AccessLog 访问日志
func AccessLog() HandlerFunc {
	return func(ctx *Context) {
		var startTime = time.Now()
		ctx.Next()
		log.Printf(
			"access: protocol=%s, path=%s, cost=%s\n",
			ctx.Writer.Protocol(),
			ctx.Request.Header.Get(XPath),
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
				log.Printf(string(msg))

				if ctx.Writer.Protocol() == ProtocolHTTP {
					_ = ctx.WriteString(http.StatusInternalServerError, "internal server error")
				}
			}
		}()
		ctx.Next()
	}
}

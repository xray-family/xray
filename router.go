package uRouter

import (
	"github.com/lxzan/uRouter/internal"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"sort"
)

type (
	// Router 路由器
	Router struct {
		// 分隔符
		separator string
		// 全局中间件
		middlewares []HandlerFunc
		// 接口映射
		routes map[string][]HandlerFunc
		// 路由列表
		routePaths map[string]uint8
		// 路由组列表
		groupPaths map[string]uint8
		// 路径匹配失败的处理
		OnNoMatch HandlerFunc
	}

	// HandlerFunc 处理函数
	HandlerFunc func(ctx *Context)
)

// New 创建路由器
func New() *Router {
	r := &Router{
		separator:   "/",
		middlewares: make([]HandlerFunc, 0),
		routes:      map[string][]HandlerFunc{},
		routePaths:  map[string]uint8{},
		groupPaths:  map[string]uint8{},
	}
	r.OnNoMatch = func(ctx *Context) {
		if ctx.Writer.Protocol() == ProtocolHTTP {
			_ = ctx.WriteString(http.StatusNotFound, "not found")
		}
	}
	return r
}

func (c *Router) checkPathConflict(m map[string]uint8, path string) {
	if _, ok := m[path]; ok {
		log.Fatalf("path=%s, msg=path conflict\n", path)
		return
	}
	m[path] = 1
}

// Use 设置全局中间件
func (c *Router) Use(middlewares ...HandlerFunc) {
	c.middlewares = append(c.middlewares, middlewares...)
}

// Group 创建路由组
func (c *Router) Group(path string, middlewares ...HandlerFunc) *Group {
	var group = &Group{
		router:      c,
		separator:   c.separator,
		path:        internal.Join1(path, c.separator),
		middlewares: append(c.middlewares, middlewares...),
	}
	c.checkPathConflict(c.groupPaths, group.path)
	return group
}

// On 监听事件
func (c *Router) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	path = internal.Join1(path, c.separator)
	var h = c.middlewares
	h = append(h, middlewares...)
	h = append(h, handler)
	c.routes[path] = h
	c.checkPathConflict(c.routePaths, path)
}

// Emit 分发事件
func (c *Router) Emit(ctx *Context) {
	path := internal.Join1(ctx.Request.Header.Get(XPath), c.separator)
	funcs, ok := c.routes[path]
	if !ok && c.OnNoMatch != nil {
		funcs = append(c.middlewares, c.OnNoMatch)
	}
	if len(funcs) == 0 {
		return
	}

	ctx.index = -1
	ctx.handlers = funcs
	ctx.Next()
}

// Display 展示接口列表
func (c *Router) Display() {
	var keys = make([]string, 0, len(c.routes))
	for k, _ := range c.routes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	log.Printf("API List")
	for _, key := range keys {
		var handlers = c.routes[key]
		var n = len(handlers)
		if n == 0 {
			continue
		}

		funcName := runtime.FuncForPC(reflect.ValueOf(handlers[n-1]).Pointer()).Name()
		log.Printf("path=%s, handler=%s", key, funcName)
	}
}

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
		// 静态路由
		staticRoutes map[string][]HandlerFunc
		// 动态路由
		dynamicRoutes *routeTree
		// 路径匹配失败的处理
		OnNoMatch HandlerFunc
	}

	// HandlerFunc 处理函数
	HandlerFunc func(ctx *Context)
)

// New 创建路由器
func New() *Router {
	r := &Router{
		separator:     "/",
		middlewares:   make([]HandlerFunc, 0),
		staticRoutes:  map[string][]HandlerFunc{},
		dynamicRoutes: newRouteTree(),
	}
	r.OnNoMatch = func(ctx *Context) {
		if ctx.Writer.Protocol() == ProtocolHTTP {
			_ = ctx.WriteString(http.StatusNotFound, "not found")
		}
	}
	return r
}

func (c *Router) showPathConflict(path string) {
	log.Panicf("path=%s, msg=path conflict\n", path)
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
	return group
}

// On 监听事件
func (c *Router) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	path = internal.Join1(path, c.separator)

	h := append(c.middlewares, middlewares...)
	h = append(h, handler)

	if !hasVar(path) {
		if _, ok := c.staticRoutes[path]; ok {
			c.showPathConflict(path)
			return
		}
		c.staticRoutes[path] = h
		return
	}

	if v := c.dynamicRoutes.Get(path); v != nil {
		c.showPathConflict(path)
		return
	}
	c.dynamicRoutes.Set(path, h)
}

// Emit 分发事件
func (c *Router) Emit(ctx *Context) {
	path := internal.Join1(ctx.Request.Header.Get(XPath), c.separator)

	// 优先匹配静态路由
	{
		funcs, ok := c.staticRoutes[path]
		if !ok && c.OnNoMatch != nil {
			funcs = append(c.middlewares, c.OnNoMatch)
		}
		if len(funcs) > 0 {
			ctx.index = -1
			ctx.handlers = funcs
			ctx.Next()
		}
		return
	}

}

// Display 展示接口列表
func (c *Router) Display() {
	var keys = make([]string, 0, len(c.staticRoutes))
	for k, _ := range c.staticRoutes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	log.Printf("API List")
	for _, key := range keys {
		var handlers = c.staticRoutes[key]
		var n = len(handlers)
		if n == 0 {
			continue
		}

		funcName := runtime.FuncForPC(reflect.ValueOf(handlers[n-1]).Pointer()).Name()
		log.Printf("path=%s, handler=%s", key, funcName)
	}
}

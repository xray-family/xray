package uRouter

import (
	"github.com/lxzan/uRouter/internal"
	"log"
	"reflect"
	"sort"
)

type (
	// Router 路由器
	Router struct {
		separator   string        // separator
		middlewares []HandlerFunc // global middlewares
		routes      map[string][]HandlerFunc
		HeaderCodec *HeaderCodec
		OnNoMatch   HandlerFunc
	}

	// HandlerFunc 处理函数
	HandlerFunc func(ctx *Context)
)

// New 创建路由器
func New() *Router {
	return &Router{
		separator:   "/",
		middlewares: make([]HandlerFunc, 0),
		routes:      map[string][]HandlerFunc{},
		HeaderCodec: TextHeader,
	}
}

// Use 设置全局中间件
func (c *Router) Use(middlewares ...HandlerFunc) {
	c.middlewares = append(c.middlewares, middlewares...)
}

// Group 创建路由组
func (c *Router) Group(prefix string, middlewares ...HandlerFunc) *Group {
	return &Group{
		router:      c,
		separator:   c.separator,
		prefix:      internal.Join1(prefix, c.separator),
		middlewares: append(c.middlewares, middlewares...),
	}
}

// On 监听事件
func (c *Router) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	path = internal.Join1(path, c.separator)
	var h = c.middlewares
	h = append(h, middlewares...)
	h = append(h, handler)
	c.routes[path] = h
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
		log.Printf("path=%s, handler=%s", key, reflect.TypeOf(handlers[n-1]).String())
	}
}

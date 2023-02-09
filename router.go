package uRouter

import (
	_ "embed"
	"github.com/lxzan/uRouter/helper"
	"github.com/lxzan/uRouter/internal"
	"net/http"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type (
	// Router 路由器
	Router struct {
		// 互斥锁, 防止有人搞骚操作, 多线程注册路由
		// mutual exclusion locks, prevent people from tampering, multi-thread registration routes
		mu *sync.Mutex

		// 分隔符
		separator string

		// 全局中间件
		// global middlewares
		chainsGlobal []HandlerFunc

		// 静态路由
		staticRoutes map[string][]HandlerFunc

		// 动态路由
		dynamicRoutes *routeTree

		// 单例模式, 确保chainsNotFound只构建一次
		// singleton pattern, ensuring that chainsNotFound is built only once
		once *sync.Once

		// 路径匹配失败的处理函数链: 全局中间件 + OnNotFound
		// path match failure handler chain: global middleware + OnNotFound
		chainsNotFound []HandlerFunc

		// 路径匹配失败的处理函数
		// path matching failure handling function
		OnNotFound HandlerFunc
	}

	// HandlerFunc 处理函数
	HandlerFunc func(ctx *Context)
)

// New 创建路由器
func New() *Router {
	r := &Router{
		separator:     defaultSeparator,
		once:          &sync.Once{},
		mu:            &sync.Mutex{},
		chainsGlobal:  make([]HandlerFunc, 0),
		staticRoutes:  map[string][]HandlerFunc{},
		dynamicRoutes: newRouteTree(),
	}
	r.OnNotFound = func(ctx *Context) {
		if ctx.Writer.Protocol() == ProtocolHTTP {
			_ = ctx.WriteString(http.StatusNotFound, "not found")
		}
	}
	return r
}

func (c *Router) showPathConflict(path string) {
	Logger().Panic("path=%s, msg=path conflict", path)
}

// pathExists 检测接口是否存在, 最坏情况O(n)复杂度
// detection of interface existence, worst-case O(n) complexity
func (c *Router) pathExists(path string) bool {
	if _, ok := c.staticRoutes[path]; ok {
		return true
	}
	if _, ok := c.dynamicRoutes.Get(path); ok {
		return true
	}

	if hasVar(path) {
		var list1 = helper.Split(path, defaultSeparator)
		var n = len(list1)
		for k, _ := range c.staticRoutes {
			var list2 = helper.Split(k, defaultSeparator)
			if n == len(list2) {
				var sum = 0
				for i, v := range list2 {
					if v == list1[i] || isVar(list1[i]) {
						sum++
					}
				}
				if sum == n {
					return true
				}
			}
		}
	}
	return false
}

// Use 设置全局中间件
// set global middlewares
func (c *Router) Use(middlewares ...HandlerFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.chainsGlobal = append(c.chainsGlobal, middlewares...)
}

// Group 创建路由组
// create a route group
func (c *Router) Group(path string, middlewares ...HandlerFunc) *Group {
	c.mu.Lock()
	defer c.mu.Unlock()

	var group = &Group{
		router:      c,
		separator:   c.separator,
		path:        internal.JoinPath(c.separator, path),
		middlewares: append(c.chainsGlobal, middlewares...),
	}
	return group
}

// On 监听事件
// listen to event
func (c *Router) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnAction("", path, handler, middlewares...)
}

// OnAction 类似On方法, 多了一个动作修饰词
func (c *Router) OnAction(action string, path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	vpath := internal.JoinPath(c.separator, strings.ToLower(action), path)
	h := append(c.chainsGlobal, middlewares...)
	h = append(h, handler)

	// 检测路径冲突
	if c.pathExists(vpath) {
		c.showPathConflict(vpath)
		return
	}

	if !hasVar(vpath) {
		c.staticRoutes[vpath] = h
	} else {
		c.dynamicRoutes.Set(vpath, h)
	}
}

// Emit 分发事件
// emit event
func (c *Router) Emit(path string, ctx *Context) {
	c.EmitAction("", path, ctx)
}

// EmitAction 分发事件
// emit event
func (c *Router) EmitAction(action string, path string, ctx *Context) {
	ctx.Request.RPath = internal.JoinPath(c.separator, strings.ToLower(action), path)

	{
		// 优先匹配静态路由
		// preferred match for static routes
		if funcs, ok := c.staticRoutes[ctx.Request.RPath]; ok {
			ctx.Request.VPath = ctx.Request.RPath
			if len(funcs) > 0 {
				ctx.handlers = funcs
				ctx.Next()
			}
			return
		}
	}

	{
		// 匹配动态路由
		// matching dynamic routes
		if h, ok := c.dynamicRoutes.Get(ctx.Request.RPath); ok {
			ctx.Request.VPath = h.VPath
			if len(h.Funcs) > 0 {
				ctx.handlers = h.Funcs
				ctx.Next()
			}
			return
		}
	}

	// 匹配失败的处理
	// handling of failed matches
	c.once.Do(func() {
		c.chainsNotFound = append(c.chainsGlobal, c.OnNotFound)
	})
	ctx.handlers = c.chainsNotFound
	ctx.Next()
}

//go:embed asserts/bless.txt
var blessMessage string

// Display 展示接口列表
// display api list
func (c *Router) Display() {
	Logger().Info(blessMessage)

	Logger().Info("uRouter is running")
	var keys = make([]string, 0, len(c.staticRoutes))
	for k, _ := range c.staticRoutes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	Logger().Info("API List:")
	for _, key := range keys {
		var handlers = c.staticRoutes[key]
		var n = len(handlers)
		if n > 0 {
			funcName := runtime.FuncForPC(reflect.ValueOf(handlers[n-1]).Pointer()).Name()
			Logger().Info("path=%s, handler=%s", key, funcName)
		}
	}

	c.dynamicRoutes.Range(func(h *apiHandler) {
		var handlers = h.Funcs
		var n = len(handlers)
		if n > 0 {
			funcName := runtime.FuncForPC(reflect.ValueOf(handlers[n-1]).Pointer()).Name()
			Logger().Info("path=%s, handler=%s", h.VPath, funcName)
		}
	})
}

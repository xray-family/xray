package uRouter

import (
	_ "embed"
	"github.com/lxzan/uRouter/internal"
	"net/http"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
)

const SEP = "/"

type (
	// Router 路由器
	Router struct {
		// 互斥锁, 防止有人搞骚操作, 多线程注册路由
		// mutual exclusion locks, prevent people from tampering, multi-thread registration routes
		mu *sync.Mutex

		// 预注册的接口
		apis []*apiHandler

		// 静态路由
		staticRoutes map[string]*apiHandler

		// 动态路由
		dynamicRoutes *routeTree

		// 全局中间件
		// global middlewares
		chainsGlobal []HandlerFunc

		// 路径匹配失败的处理函数链: chainsGlobal + OnNotFound
		// path match failure handler chain: chainsGlobal + OnNotFound
		chainsNotFound []HandlerFunc

		// 路径匹配失败的处理函数
		// path matching failure handling function
		OnNotFound HandlerFunc
	}

	apiHandler struct {
		Action   string        // 动作修饰词
		Path     string        // 接口路径
		FullPath string        // 接口路径, 包含action和name
		Funcs    []HandlerFunc // 处理链
	}

	// HandlerFunc 处理函数
	HandlerFunc func(ctx *Context)
)

// New 创建路由器
func New() *Router {
	r := &Router{
		apis:          []*apiHandler{},
		mu:            &sync.Mutex{},
		chainsGlobal:  make([]HandlerFunc, 0),
		staticRoutes:  map[string]*apiHandler{},
		dynamicRoutes: newRouteTree(),
	}
	r.OnNotFound = func(ctx *Context) {
		if ctx.Writer.Protocol() == ProtocolHTTP {
			_ = ctx.WriteString(http.StatusNotFound, "not found")
		}
	}
	return r
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
		path:        internal.JoinPath(SEP, path),
		middlewares: append(c.chainsGlobal, middlewares...),
	}
	return group
}

// On 监听事件
// listen to event
func (c *Router) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent("", path, handler, middlewares...)
}

// OnEvent 类似On方法, 多了一个动作修饰词
func (c *Router) OnEvent(action string, path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	action = strings.ToLower(action)
	h := append(c.chainsGlobal, middlewares...)
	h = append(h, handler)
	c.apis = append(c.apis, &apiHandler{
		Action:   action,
		Path:     internal.JoinPath(SEP, path),
		FullPath: internal.JoinPath(SEP, action, path),
		Funcs:    h,
	})
}

// Emit 分发事件
// emit event
func (c *Router) Emit(path string, ctx *Context) {
	c.EmitEvent("", path, ctx)
}

// EmitEvent 分发事件
// emit event
func (c *Router) EmitEvent(action string, path string, ctx *Context) {
	action = strings.ToLower(action)
	path = internal.JoinPath(SEP, path)
	ctx.Request.RPath = path
	fullpath := internal.JoinPath(SEP, SEP+action+path)

	{
		// 优先匹配静态路由
		// preferred match for static routes
		if h, ok := c.staticRoutes[fullpath]; ok {
			ctx.Request.VPath = h.Path
			if len(h.Funcs) > 0 {
				ctx.handlers = h.Funcs
				ctx.Next()
			}
			return
		}
	}

	{
		// 匹配动态路由
		// matching dynamic routes
		if h, ok := c.dynamicRoutes.Get(fullpath); ok {
			ctx.Request.VPath = h.Path
			if len(h.Funcs) > 0 {
				ctx.handlers = h.Funcs
				ctx.Next()
			}
			return
		}
	}

	// 匹配失败的处理
	// handling of failed matches
	ctx.handlers = c.chainsNotFound
	ctx.Next()
}

//go:embed asserts/bless.txt
var blessMessage string

// display 展示接口列表
// display api list
func (c *Router) display() {
	sort.Slice(c.apis, func(i, j int) bool {
		a := c.apis[i]
		b := c.apis[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		return a.Action < b.Action
	})

	Logger().Info(blessMessage + "\n\n")
	Logger().Info("uRouter is running")
	Logger().Info("API List:")
	for _, v := range c.apis {
		n := len(v.Funcs)
		funcName := runtime.FuncForPC(reflect.ValueOf(v.Funcs[n-1]).Pointer()).Name()
		Logger().Info("action=%s, path=%s, handler=%s", v.Action, v.Path, funcName)
	}
}

// Start 启动路由器
// 打印问候语和API列表
func (c *Router) Start() {
	c.doStart()
	c.display()
}

// 启动路由器
func (c *Router) doStart() {
	c.chainsNotFound = append(c.chainsGlobal, c.OnNotFound)

	var staticAPIs []*apiHandler
	var dynamicAPIs []*apiHandler
	for i, v := range c.apis {
		if hasVar(v.Path) {
			dynamicAPIs = append(dynamicAPIs, c.apis[i])
		} else {
			staticAPIs = append(staticAPIs, c.apis[i])
		}
	}

	var logger = func(v *apiHandler) {
		Logger().Panic("action=%s, path=%s, msg=api path conflict", v.Action, v.Path)
	}

	for i, v := range dynamicAPIs {
		if _, exist := c.dynamicRoutes.Get(v.FullPath); exist {
			logger(v)
			return
		}
		c.dynamicRoutes.Set(dynamicAPIs[i])
	}

	for i, v := range staticAPIs {
		if _, exist := c.dynamicRoutes.Get(v.FullPath); exist {
			logger(v)
			return
		}
		if _, exist := c.staticRoutes[v.FullPath]; exist {
			logger(v)
			return
		}
		c.staticRoutes[v.FullPath] = staticAPIs[i]
	}
}

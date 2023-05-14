package uRouter

import (
	_ "embed"
	"github.com/lxzan/uRouter/internal"
	"net/http"
	"reflect"
	"runtime"
	"sort"
	"time"
)

const (
	SEP     = internal.Separator
	UPath   = "U-Path"
	UAction = "U-Action"
)

type (
	// Router 路由器
	Router struct {
		// 静态路由
		staticMatcher staticMatcher

		// 动态路由
		dynamicMatcher dynamicMatcher

		// 全局中间件
		// global middlewares
		chainsGlobal []HandlerFunc

		// 路径匹配失败的处理函数链: chainsGlobal + onNotFound
		// path match failure handler chain: chainsGlobal + onNotFound
		chainsNotFound []HandlerFunc

		// 路径匹配失败的处理函数
		// path matching failure handling function
		onNotFound HandlerFunc
	}

	apiHandler struct {
		Action string        // 动作修饰词
		Path   string        // 接口路径
		Funcs  []HandlerFunc // 处理链
	}

	// HandlerFunc 处理函数
	HandlerFunc func(ctx *Context)
)

// New 创建路由器
func New() *Router {
	r := &Router{
		chainsGlobal:   make([]HandlerFunc, 0),
		staticMatcher:  make(staticMatcher, 0),
		dynamicMatcher: make(dynamicMatcher),
	}

	r.onNotFound = func(ctx *Context) {
		if ctx.Writer.Protocol() == ProtocolHTTP {
			_ = ctx.WriteString(http.StatusNotFound, "not found")
		}
	}

	go func() {
		time.Sleep(time.Second)
		r.display()
	}()

	return r
}

// Use 设置全局中间件
// set global middlewares
func (c *Router) Use(middlewares ...HandlerFunc) {
	c.chainsGlobal = append(c.chainsGlobal, middlewares...)
	c.SetHandlerNotFound(c.onNotFound)
}

// SetHandlerNotFound 设置找不到路由的处理函数
func (c *Router) SetHandlerNotFound(handler HandlerFunc) {
	c.onNotFound = handler
	c.chainsNotFound = internal.Clone(c.chainsGlobal)
	c.chainsNotFound = append(c.chainsNotFound, handler)
}

// Group 创建路由组
// create a route group
func (c *Router) Group(path string, middlewares ...HandlerFunc) *Group {
	var group = &Group{
		router:      c,
		path:        internal.JoinPath(SEP, path),
		middlewares: append(internal.Clone(c.chainsGlobal), middlewares...),
	}
	return group
}

// On 监听事件
// listen to event
func (c *Router) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent("", path, handler, middlewares...)
}

func (c *Router) OnGET(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent(http.MethodGet, path, handler, middlewares...)
}

func (c *Router) OnPOST(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent(http.MethodPost, path, handler, middlewares...)
}

func (c *Router) OnPUT(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent(http.MethodPut, path, handler, middlewares...)
}

func (c *Router) OnDELETE(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent(http.MethodDelete, path, handler, middlewares...)
}

// 报告路由冲突
func (c *Router) reportConflict(api1, api2 *apiHandler) {
	Logger().Panic("action=%s, path=[ %s, %s ], msg=api path conflict", api1.Action, api1.Path, api2.Path)
}

// OnEvent 监听一个事件, 绑定处理函数
// action: 操作修饰词, 区分大小写, 可为空字符串; 在HTTP里固定为Method
// path: 请求路径
// handler: 处理函数
// listens for an event, binds a handler
// action: action modifier, case-sensitive, can be an empty string; fixed to Method in HTTP
// path: request path
// handler: handler function
func (c *Router) OnEvent(action string, path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	h := append(internal.Clone(c.chainsGlobal), middlewares...)
	h = append(h, handler)
	api := &apiHandler{
		Action: action,
		Path:   internal.JoinPath(SEP, path),
		Funcs:  h,
	}
	setApiHandler(c, api.Action, api.Path, api)
}

// Emit 分发事件
// emit event
func (c *Router) Emit(path string, ctx *Context) {
	c.EmitEvent("", path, ctx)
}

// EmitEvent 分发事件
// emit event
func (c *Router) EmitEvent(action string, path string, ctx *Context) {
	path = internal.TrimPath(path)
	ctx.Request.RPath = path

	if h, ok := getApiHandler(c, action, path); ok {
		ctx.Request.VPath = h.Path
		ctx.handlers = h.Funcs
		ctx.Next()
		return
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
	var apis []*apiHandler
	c.staticMatcher.Range(func(h *apiHandler) { apis = append(apis, h) })
	c.dynamicMatcher.Range(func(h *apiHandler) { apis = append(apis, h) })
	sort.Slice(apis, func(i, j int) bool {
		a := apis[i]
		b := apis[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		return a.Action < b.Action
	})

	Logger().Info(blessMessage + "\n\n")
	Logger().Info("uRouter is running")
	Logger().Info("API List:")

	var actions []string
	var paths []string
	for _, v := range apis {
		actions = append(actions, v.Action)
		paths = append(paths, v.Path)
	}
	actionLength := internal.GetMaxLength(actions...)
	pathLength := internal.GetMaxLength(paths...)

	for _, v := range apis {
		n := len(v.Funcs)
		funcName := runtime.FuncForPC(reflect.ValueOf(v.Funcs[n-1]).Pointer()).Name()
		Logger().Info(
			"action=%s path=%s handler=%s",
			internal.Padding(v.Action, actionLength),
			internal.Padding(v.Path, pathLength),
			funcName,
		)
	}
}

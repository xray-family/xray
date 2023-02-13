package uRouter

import (
	_ "embed"
	"github.com/lxzan/uRouter/internal"
	"net/http"
	"reflect"
	"runtime"
	"sort"
)

const (
	SEP     = internal.Separator
	UPath   = "U-Path"
	UAction = "U-Action"
)

type (
	// Router 路由器
	Router struct {
		// 预注册的接口
		apis []*apiHandler

		// 静态路由
		staticRoutes map[string]map[string]*apiHandler

		// 动态路由
		dynamicRoutes map[string]*routeTree

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
		apis:          []*apiHandler{},
		chainsGlobal:  make([]HandlerFunc, 0),
		staticRoutes:  map[string]map[string]*apiHandler{},
		dynamicRoutes: map[string]*routeTree{},
	}
	r.OnNotFound = func(ctx *Context) {
		if ctx.Writer.Protocol() == ProtocolHTTP {
			_ = ctx.WriteString(http.StatusNotFound, "not found")
		}
	}
	return r
}

// cloneMiddlewares deep clone
func (c *Router) cloneMiddlewares(chains []HandlerFunc) []HandlerFunc {
	var results []HandlerFunc
	for i, _ := range chains {
		results = append(results, chains[i])
	}
	return results
}

// Use 设置全局中间件
// set global middlewares
func (c *Router) Use(middlewares ...HandlerFunc) {
	c.chainsGlobal = append(c.chainsGlobal, middlewares...)
}

// Group 创建路由组
// create a route group
func (c *Router) Group(path string, middlewares ...HandlerFunc) *Group {
	var group = &Group{
		router:      c,
		path:        internal.JoinPath(SEP, path),
		middlewares: append(c.cloneMiddlewares(c.chainsGlobal), middlewares...),
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

// OnEvent 监听一个事件, 绑定处理函数
// action: 操作修饰词, 区分大小写, 可为空字符串; 在HTTP里固定为Method
// path: 请求路径
// handler: 处理函数
// listens for an event, binds a handler
// action: action modifier, case-sensitive, can be an empty string; fixed to Method in HTTP
// path: request path
// handler: handler function
func (c *Router) OnEvent(action string, path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	h := append(c.cloneMiddlewares(c.chainsGlobal), middlewares...)
	h = append(h, handler)
	c.apis = append(c.apis, &apiHandler{
		Action: action,
		Path:   internal.JoinPath(SEP, path),
		Funcs:  h,
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

	var actions []string
	var paths []string
	for _, v := range c.apis {
		actions = append(actions, v.Action)
		paths = append(paths, v.Path)
	}
	actionLength := internal.GetMaxLength(actions...)
	pathLength := internal.GetMaxLength(paths...)

	for _, v := range c.apis {
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

// Start
// 启动路由器并打印问候语和API列表
func (c *Router) Start() {
	c.StartSilently()
	c.display()
}

// StartSilently 仅启动路由器, 不打印问候语和API列表
func (c *Router) StartSilently() {
	c.chainsNotFound = append(c.cloneMiddlewares(c.chainsGlobal), c.OnNotFound)

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
		if _, exist := getApiHandler(c, v.Action, v.Path); exist {
			logger(v)
			return
		}
		setApiHandler(c, v.Action, v.Path, dynamicAPIs[i])
	}

	for i, v := range staticAPIs {
		if _, exist := getApiHandler(c, v.Action, v.Path); exist {
			logger(v)
			return
		}
		setApiHandler(c, v.Action, v.Path, staticAPIs[i])
	}
}

package xray

import (
	_ "embed"
	"github.com/xray-family/xray/codec"
	"github.com/xray-family/xray/internal"
	"github.com/xray-family/xray/internal/treemap"
	"github.com/xray-family/xray/log"
	"net/http"
	"reflect"
	"runtime"
	"sort"
	"time"
)

const _sep = internal.Separator

type (
	// Router 路由器
	Router struct {
		// 配置
		conf *config

		// 路由树
		matcher *treemap.TreeMap[*apiHandler]

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
		Method string        // 操作方法
		Path   string        // 接口路径
		Funcs  []HandlerFunc // 处理链
	}

	// HandlerFunc 处理函数
	HandlerFunc func(ctx *Context)
)

// New 创建路由器
func New(options ...Option) *Router {
	conf := &config{greeting: greeting{enabled: true, delay: time.Second}}
	options = append(options, withInit())
	for _, f := range options {
		f(conf)
	}

	r := &Router{
		conf:         conf,
		chainsGlobal: make([]HandlerFunc, 0),
		matcher:      treemap.New[*apiHandler](),
	}

	r.SetHandlerNotFound(func(ctx *Context) {
		if ctx.Writer.Protocol() == ProtocolHTTP {
			_ = ctx.WriteString(http.StatusNotFound, "not found")
		}
	})

	if conf.greeting.enabled {
		go func() {
			time.Sleep(conf.greeting.delay)
			r.display()
		}()
	}

	return r
}

// JsonCodec 获取JSON编码器
func (c *Router) JsonCodec() codec.Codec {
	return c.conf.jsonCodec
}

// Logger 获取日志工具
func (c *Router) Logger() log.Logger {
	return c.conf.logger
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
		path:        internal.JoinPath(_sep, path),
		middlewares: append(internal.Clone(c.chainsGlobal), middlewares...),
	}
	return group
}

// On 监听事件
// listen to event
func (c *Router) On(path string, handlers ...HandlerFunc) {
	c.OnEvent("", path, handlers...)
}

func (c *Router) GET(path string, handlers ...HandlerFunc) {
	c.OnEvent(http.MethodGet, path, handlers...)
}

func (c *Router) POST(path string, handlers ...HandlerFunc) {
	c.OnEvent(http.MethodPost, path, handlers...)
}

func (c *Router) PUT(path string, handlers ...HandlerFunc) {
	c.OnEvent(http.MethodPut, path, handlers...)
}

func (c *Router) DELETE(path string, handlers ...HandlerFunc) {
	c.OnEvent(http.MethodDelete, path, handlers...)
}

// 报告路由冲突
func (c *Router) reportConflict(api1, api2 *apiHandler) {
	c.Logger().Panic("method=%s, path=[ %s, %s ], msg=path conflict", api1.Method, api1.Path, api2.Path)
}

// OnEvent 监听一个事件, 绑定处理函数
// method: 操作修饰词, 区分大小写, 可以为空
// path: 请求路径
// handler: 处理函数
// listens for an event, binds a handler
// method: action modifier, case-sensitive, can be an empty string;
// path: request path
// handler: handler function
func (c *Router) OnEvent(method string, path string, handlers ...HandlerFunc) {
	api := &apiHandler{
		Method: method,
		Path:   internal.JoinPath(_sep, path),
		Funcs:  append(internal.Clone(c.chainsGlobal), handlers...),
	}
	if v, exists := c.matcher.Exists(api.Method, api.Path); exists {
		c.reportConflict(api, v)
	}
	c.matcher.Set(api.Method, api.Path, api)
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

	if h, ok := c.matcher.Get(action, path); ok {
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
	c.matcher.Range(func(h *apiHandler) { apis = append(apis, h) })
	sort.Slice(apis, func(i, j int) bool {
		a := apis[i]
		b := apis[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		return a.Method < b.Method
	})

	c.conf.logger.Info(blessMessage + "\n\n")
	c.conf.logger.Info("Xray is running")
	c.conf.logger.Info("API List:")

	var actions []string
	var paths []string
	for _, v := range apis {
		actions = append(actions, v.Method)
		paths = append(paths, v.Path)
	}
	actionLength := internal.GetMaxLength(actions...)
	pathLength := internal.GetMaxLength(paths...)

	for _, v := range apis {
		n := len(v.Funcs)
		funcName := runtime.FuncForPC(reflect.ValueOf(v.Funcs[n-1]).Pointer()).Name()
		c.conf.logger.Info(
			"action=%s path=%s handler=%s",
			internal.Padding(v.Method, actionLength),
			internal.Padding(v.Path, pathLength),
			funcName,
		)
	}
}

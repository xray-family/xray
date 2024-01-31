package xray

import (
	"github.com/xray-family/xray/internal"
	"net/http"
)

// Group 路由组
// route group
type Group struct {
	router      *Router
	path        string
	middlewares []HandlerFunc
}

// Group 创建子路由组
// create a child route group
func (c *Group) Group(path string, middlewares ...HandlerFunc) *Group {
	group := &Group{
		router:      c.router,
		path:        internal.JoinPath(_sep, c.path, path),
		middlewares: append(internal.Clone(c.middlewares), middlewares...),
	}
	return group
}

// OnEvent 监听事件
// listen to event
func (c *Group) OnEvent(method string, path string, handlers ...HandlerFunc) {
	h := append(internal.Clone(c.middlewares), handlers...)
	api := &apiHandler{
		Method: method,
		Path:   internal.JoinPath(_sep, c.path, path),
		Funcs:  h,
	}
	setApiHandler(c.router, api.Method, api.Path, api)
}

// On  类似OnEvent方法, 但是没有操作方法
func (c *Group) On(path string, handlers ...HandlerFunc) {
	c.OnEvent("", path, handlers...)
}

func (c *Group) GET(path string, handlers ...HandlerFunc) {
	c.OnEvent(http.MethodGet, path, handlers...)
}

func (c *Group) POST(path string, handlers ...HandlerFunc) {
	c.OnEvent(http.MethodPost, path, handlers...)
}

func (c *Group) PUT(path string, handlers ...HandlerFunc) {
	c.OnEvent(http.MethodPut, path, handlers...)
}

func (c *Group) DELETE(path string, handlers ...HandlerFunc) {
	c.OnEvent(http.MethodDelete, path, handlers...)
}

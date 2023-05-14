package uRouter

import (
	"github.com/lxzan/uRouter/internal"
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
		path:        internal.JoinPath(SEP, c.path, path),
		middlewares: append(internal.Clone(c.middlewares), middlewares...),
	}
	return group
}

// OnEvent 监听事件
// listen to event
func (c *Group) OnEvent(action string, path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	h := append(internal.Clone(c.middlewares), middlewares...)
	h = append(h, handler)
	api := &apiHandler{
		Action: action,
		Path:   internal.JoinPath(SEP, c.path, path),
		Funcs:  h,
	}
	setApiHandler(c.router, api.Action, api.Path, api)
}

// On  类似OnEvent方法, 但是没有动作修饰词
func (c *Group) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent("", path, handler, middlewares...)
}

func (c *Group) OnGET(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent(http.MethodGet, path, handler, middlewares...)
}

func (c *Group) OnPOST(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent(http.MethodPost, path, handler, middlewares...)
}

func (c *Group) OnPUT(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent(http.MethodPut, path, handler, middlewares...)
}

func (c *Group) OnDELETE(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent(http.MethodDelete, path, handler, middlewares...)
}

package uRouter

import "github.com/lxzan/uRouter/internal"

// Group 路由组
type Group struct {
	router      *Router
	separator   string // separator
	path        string
	middlewares []HandlerFunc
}

// Group 创建子路由组
func (c *Group) Group(path string, middlewares ...HandlerFunc) *Group {
	group := &Group{
		router:      c.router,
		separator:   c.separator,
		path:        internal.Join2(c.path, path, c.separator),
		middlewares: append(c.middlewares, middlewares...),
	}
	c.router.checkPathConflict(c.router.groupPaths, group.path)
	return group
}

// On 监听事件
func (c *Group) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	path = internal.Join2(c.path, path, c.separator)
	var h = c.middlewares
	h = append(h, middlewares...)
	h = append(h, handler)
	c.router.routes[path] = h
	c.router.checkPathConflict(c.router.routePaths, path)
}

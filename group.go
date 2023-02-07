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
	return group
}

// On 监听事件
func (c *Group) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	path = internal.Join2(c.path, path, c.separator)
	router := c.router
	h := append(c.middlewares, middlewares...)
	h = append(h, handler)

	if !hasVar(path) {
		if _, ok := router.staticRoutes[path]; ok {
			router.showPathConflict(path)
			return
		}
		router.staticRoutes[path] = h
		return
	}

	if v := router.dynamicRoutes.Get(path); v != nil {
		router.showPathConflict(path)
		return
	}
	router.dynamicRoutes.Set(path, h)
}

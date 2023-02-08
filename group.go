package uRouter

import "github.com/lxzan/uRouter/internal"

// Group 路由组
// route group
type Group struct {
	router      *Router
	separator   string
	path        string
	middlewares []HandlerFunc
}

// Group 创建子路由组
// create a child route group
func (c *Group) Group(path string, middlewares ...HandlerFunc) *Group {
	c.router.mu.Lock()
	defer c.router.mu.Unlock()

	group := &Group{
		router:      c.router,
		separator:   c.separator,
		path:        internal.Join2(c.path, path, c.separator),
		middlewares: append(c.middlewares, middlewares...),
	}
	return group
}

// On 监听事件
// listen to event
func (c *Group) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	router := c.router
	router.mu.Lock()
	defer router.mu.Unlock()

	path = internal.Join2(c.path, path, c.separator)
	h := append(c.middlewares, middlewares...)
	h = append(h, handler)

	// 检测路径冲突
	if router.pathExists(path) {
		router.showPathConflict(path)
		return
	}

	if !hasVar(path) {
		router.staticRoutes[path] = h
	} else {
		router.dynamicRoutes.Set(path, h)
	}
}

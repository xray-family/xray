package uRouter

// Group 路由组
type Group struct {
	router      *Router
	separator   string // separator
	prefix      string
	middlewares []HandlerFunc
}

// Group 创建子路由组
func (c *Group) Group(prefix string, middlewares ...HandlerFunc) *Group {
	return &Group{
		router:      c.router,
		separator:   c.separator,
		prefix:      join2(c.prefix, prefix, c.separator),
		middlewares: append(c.middlewares, middlewares...),
	}
}

// On 监听事件
func (c *Group) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	path = join2(c.prefix, path, c.separator)
	var h = c.middlewares
	h = append(h, middlewares...)
	h = append(h, handler)
	c.router.routes[path] = h
}

package xray

import (
	"github.com/lxzan/xray/internal"
)

const defaultVarPrefix = ':' // 默认变量前缀

type matcher interface {
	Set(handler *apiHandler)
	Get(action, path string) (*apiHandler, bool)
	Range(f func(h *apiHandler))
}

type routeTree struct {
	Value    *apiHandler
	Children map[string]*routeTree
}

func newRouteTree() *routeTree {
	return &routeTree{
		Children: make(map[string]*routeTree),
	}
}

// 判断字符串是否为变量
func isVar(s string) bool {
	return len(s) > 0 && s[0] == defaultVarPrefix
}

// 判断字符串是否含有变量
func hasVar(s string) bool {
	var n = len(s)
	for i := 0; i < n-1; i++ {
		if s[i] == _sep[0] && s[i+1] == defaultVarPrefix {
			return true
		}
	}
	return false
}

func (c *routeTree) Set(handler *apiHandler) {
	var list = internal.Split(handler.Path)
	if len(list) == 0 {
		return
	}
	c.doSet(c, 0, list, handler)
}

func (c *routeTree) doSet(node *routeTree, index int, list []string, handler *apiHandler) {
	var segment = list[index]
	if isVar(segment) {
		segment = "*"
	}

	if node.Children == nil {
		node.Children = make(map[string]*routeTree)
	}

	var next = node.Children[segment]
	if node.Children[segment] == nil {
		next = &routeTree{}
		node.Children[segment] = next
	}

	if index+1 == len(list) {
		next.Value = handler
	} else {
		c.doSet(next, index+1, list, handler)
	}
}

func (c *routeTree) Get(path string) (*apiHandler, bool) {
	var tree = c
	var expected = 0
	var actual = 0

	internal.FastSplit(path, func(str string) bool {
		expected++
		if v, ok := tree.Children[str]; ok {
			tree = v
			actual++
			return true
		}
		if v, ok := tree.Children["*"]; ok {
			tree = v
			actual++
			return true
		}
		return false
	})

	if expected == actual {
		return tree.Value, tree.Value != nil
	}
	return nil, false
}

func (c *routeTree) Range(f func(h *apiHandler)) {
	c.doRange(c, f)
}

func (c *routeTree) doRange(node *routeTree, f func(h *apiHandler)) {
	if node == nil {
		return
	}
	if node.Value != nil {
		f(node.Value)
	}
	for _, v := range node.Children {
		c.doRange(v, f)
	}
}

type staticMatcher map[string]map[string]*apiHandler

func (c staticMatcher) Set(h *apiHandler) {
	if _, ok := c[h.Method]; !ok {
		c[h.Method] = make(map[string]*apiHandler)
	}
	c[h.Method][h.Path] = h
}

func (c staticMatcher) Get(action, path string) (*apiHandler, bool) {
	v1, ok1 := c[action]
	if !ok1 {
		return nil, false
	}
	v2, ok2 := v1[path]
	return v2, ok2
}

func (c staticMatcher) Range(f func(h *apiHandler)) {
	for _, routes := range c {
		for _, handler := range routes {
			f(handler)
		}
	}
}

type dynamicMatcher map[string]*routeTree

func (c dynamicMatcher) Set(h *apiHandler) {
	if _, ok := c[h.Method]; !ok {
		c[h.Method] = newRouteTree()
	}
	c[h.Method].Set(h)
}

func (c dynamicMatcher) Get(action, path string) (*apiHandler, bool) {
	v1, ok1 := c[action]
	if !ok1 {
		return nil, false
	}
	return v1.Get(path)
}

func (c dynamicMatcher) Range(f func(h *apiHandler)) {
	for _, tree := range c {
		tree.Range(f)
	}
}

// 查找接口绑定的处理函数
func getApiHandler(r *Router, action string, path string) (*apiHandler, bool) {
	if v, ok := r.staticMatcher.Get(action, path); ok {
		return v, ok
	}
	return r.dynamicMatcher.Get(action, path)
}

// 设置接口绑定的处理函数
func setApiHandler(r *Router, action string, path string, api *apiHandler) {
	// 检测冲突
	if v, exists := getApiHandler(r, action, path); exists {
		r.reportConflict(api, v)
	}

	// 静态路由
	if !hasVar(path) {
		r.staticMatcher.Set(api)
		return
	}

	// 动态路由
	r.dynamicMatcher.Set(api)

	// 检测冲突
	r.staticMatcher.Range(func(h *apiHandler) {
		if _, exists := r.dynamicMatcher.Get(h.Method, h.Path); exists {
			r.reportConflict(api, h)
		}
	})
}

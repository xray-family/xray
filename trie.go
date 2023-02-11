package uRouter

import (
	"github.com/lxzan/uRouter/internal"
)

const defaultVarPrefix = ':' // 默认变量前缀

type (
	routeTree struct {
		Value    *apiHandler
		Children map[string]*routeTree
	}
)

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
		if s[i] == SEP[0] && s[i+1] == defaultVarPrefix {
			return true
		}
	}
	return false
}

func (c *routeTree) Set(handler *apiHandler) {
	var list = internal.Split(handler.FullPath, SEP)
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
	var list = internal.Split(path, SEP)
	if len(list) == 0 {
		return nil, false
	}
	return c.doGet(c, 0, list)
}

func (c *routeTree) doGet(node *routeTree, index int, list []string) (*apiHandler, bool) {
	if index == len(list) {
		return node.Value, node.Value != nil
	}
	var segment = list[index]
	if v, ok := node.Children[segment]; ok {
		return c.doGet(v, index+1, list)
	}
	if v, ok := node.Children["*"]; ok {
		return c.doGet(v, index+1, list)
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

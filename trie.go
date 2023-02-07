package uRouter

import "github.com/lxzan/uRouter/internal"

const (
	defaultSeparator = "/"
	defaultVarPrefix = ':'
)

type (
	apiHandler struct {
		VPath string
		Funcs []HandlerFunc
	}

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

func isVar(s string) bool {
	return len(s) > 0 && s[0] == defaultVarPrefix
}

func hasVar(s string) bool {
	var list = internal.Split(s, defaultSeparator)
	for _, item := range list {
		if isVar(item) {
			return true
		}
	}
	return false
}

func (c *routeTree) Set(vpath string, handlers []HandlerFunc) {
	var handler = &apiHandler{VPath: vpath, Funcs: handlers}
	var list = internal.Split(handler.VPath, defaultSeparator)
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

func (c *routeTree) Get(path string) *apiHandler {
	var list = internal.Split(path, defaultSeparator)
	if len(list) == 0 {
		return nil
	}
	return c.doGet(c, 0, list)
}

func (c *routeTree) doGet(node *routeTree, index int, list []string) *apiHandler {
	if index == len(list) {
		return node.Value
	}
	var segment = list[index]
	if v, ok := node.Children[segment]; ok {
		return c.doGet(v, index+1, list)
	}
	if v, ok := node.Children["*"]; ok {
		return c.doGet(v, index+1, list)
	}
	return nil
}

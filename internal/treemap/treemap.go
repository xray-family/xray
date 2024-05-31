package treemap

import (
	"strings"
)

const (
	SEP      = "/"
	Wildcard = "*"
)

func New[T any]() *TreeMap[T] {
	return &TreeMap[T]{
		static: make(map[string]map[string]*Element[T]),
		root:   make(map[string]*Element[T]),
	}
}

type (
	TreeMap[T any] struct {
		static map[string]map[string]*Element[T]
		root   map[string]*Element[T]
	}

	Element[T any] struct {
		children map[string]*Element[T]
		ok       bool
		key      string
		value    T
	}
)

func (c *TreeMap[T]) trimLeft(path string) string {
	for n := len(path); n > 0 && path[0] == SEP[0]; n-- {
		path = path[1:]
	}
	return path
}

// 判断字符串是否含有变量
func (c *TreeMap[T]) hasVar(s string) bool { return strings.Contains(s, "/:") }

func (c *TreeMap[T]) formatKey(key string) string {
	if len(key) > 0 && key[0] == ':' {
		return Wildcard
	}
	return key
}

func (c *TreeMap[T]) Set(method, key string, val T) {
	if c.static[method] == nil {
		c.static[method] = make(map[string]*Element[T])
	}
	if c.root[method] == nil {
		c.root[method] = new(Element[T])
	}
	var node = &Element[T]{key: key, value: val, ok: true}
	c.doSet(c.root[method], node, key)
	if !c.hasVar(key) {
		c.static[method][key] = node
	}
}

func (c *TreeMap[T]) doSet(far, son *Element[T], key string) {
	key = c.trimLeft(key)
	if far.children == nil {
		far.children = make(map[string]*Element[T])
	}

	index := strings.Index(key, SEP)
	if index < 0 {
		key = c.formatKey(key)
		if node := far.children[key]; node == nil {
			far.children[key] = son
		} else {
			node.value, node.ok = son.value, son.ok
		}
		return
	}

	s0, s1 := c.formatKey(key[:index]), key[index+1:]
	next := far.children[s0]
	if far.children[s0] == nil {
		next = &Element[T]{}
		far.children[s0] = next
	}
	c.doSet(next, son, s1)
}

func (c *TreeMap[T]) Get(method, key string) (value T, exist bool) {
	if node := c.static[method]; node != nil {
		if t, ok := node[key]; ok {
			return t.value, true
		}
	}
	if node := c.root[method]; node != nil {
		c.doGet(node, key, func(t *Element[T]) { value, exist = t.value, true })
	}
	return value, exist
}

func (c *TreeMap[T]) doGet(cur *Element[T], key string, cb func(*Element[T])) {
	key = c.trimLeft(key)
	index := strings.Index(key, SEP)
	if index < 0 {
		if v, ok := cur.children[key]; ok && v.ok {
			cb(v)
		}
		if v, ok := cur.children[Wildcard]; ok && v.ok {
			cb(v)
		}
		return
	}

	s0, s1 := key[:index], key[index+1:]
	if v, ok := cur.children[s0]; ok {
		c.doGet(v, s1, cb)
	}
	if v, ok := cur.children[Wildcard]; ok {
		c.doGet(v, s1, cb)
	}
}

// Exists 检查路由冲突
func (c *TreeMap[T]) Exists(method, key string) (value T, exist bool) {
	if node := c.static[method]; node != nil {
		if t, ok := node[key]; ok {
			return t.value, true
		}
	}
	if node := c.root[method]; node != nil {
		c.doExists(node, key, func(t *Element[T]) { value, exist = t.value, true })
	}
	return value, exist
}

func (c *TreeMap[T]) doExists(cur *Element[T], key string, cb func(*Element[T])) {
	key = c.trimLeft(key)
	index := strings.Index(key, SEP)
	if index < 0 {
		key = c.formatKey(key)
		if key == Wildcard {
			for _, v := range cur.children {
				if v.ok {
					cb(v)
				}
			}
		} else {
			if v, ok := cur.children[key]; ok && v.ok {
				cb(v)
			}
			if v, ok := cur.children[Wildcard]; ok && v.ok {
				cb(v)
			}
		}
		return
	}

	s0, s1 := c.formatKey(key[:index]), key[index+1:]
	if s0 == Wildcard {
		for _, v := range cur.children {
			c.doExists(v, s1, cb)
		}
	} else {
		if v, ok := cur.children[s0]; ok {
			c.doExists(v, s1, cb)
		}
		if v, ok := cur.children[Wildcard]; ok {
			c.doExists(v, s1, cb)
		}
	}
}

func (c *TreeMap[T]) Range(f func(h T)) {
	for _, v := range c.root {
		c.doRange(v, f)
	}
}

func (c *TreeMap[T]) doRange(node *Element[T], f func(h T)) {
	if node == nil {
		return
	}
	if node.ok {
		f(node.value)
	}
	for _, v := range node.children {
		c.doRange(v, f)
	}
}

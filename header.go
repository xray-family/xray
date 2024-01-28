package xray

import (
	"net/http"
)

type Header interface {
	New() Header                          // 创建Header实例
	Set(key, value string)                // 设置键值对
	Get(key string) string                // 获取一个值
	Del(key string)                       // 删除
	Len() int                             // 获取长度
	Range(f func(key, value string) bool) // 遍历
}

type HttpHeader struct {
	http.Header
}

func (c HttpHeader) New() Header { return nil }

func (c HttpHeader) Len() int {
	return len(c.Header)
}

func (c HttpHeader) Range(f func(key string, value string) bool) {
	for k, values := range c.Header {
		for _, v := range values {
			f(k, v)
		}
	}
}

type SliceHeader [][2]string

func (c *SliceHeader) New() Header {
	return &SliceHeader{}
}

func (c *SliceHeader) Set(key, value string) {
	vec := c.Elem()
	for i, _ := range vec {
		if vec[i][0] == key {
			vec[i][1] = value
			return
		}
	}
	*c = append(*c, [2]string{key, value})
}

func (c *SliceHeader) Get(key string) (val string) {
	for _, item := range c.Elem() {
		if item[0] == key {
			return item[1]
		}
	}
	return val
}

func (c *SliceHeader) Del(key string) {
	var n = c.Len()
	for i, item := range c.Elem() {
		if item[0] == key {
			(*c)[i], (*c)[n-1] = (*c)[n-1], (*c)[i]
			*c = (*c)[:n-1]
			return
		}
	}
}

func (c *SliceHeader) Len() int {
	return len(*c)
}

func (c *SliceHeader) Elem() [][2]string {
	return *c
}

func (c *SliceHeader) Range(f func(key string, value string) bool) {
	for _, item := range c.Elem() {
		if !f(item[0], item[1]) {
			return
		}
	}
}

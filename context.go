package xray

import (
	"github.com/lxzan/xray/constant"
	"github.com/lxzan/xray/internal"
	"io"
)

type (
	// ResponseWriter 响应写入器
	ResponseWriter interface {
		// Protocol 获取协议
		Protocol() string

		// Header 获取header
		Header() Header

		// Write 写入数据
		Write([]byte) (int, error)

		// Code 设置状态码
		Code(code int)

		// Flush 刷新缓冲区
		Flush() error

		// Raw 获取原生响应写入器
		Raw() any
	}

	// Request 请求
	Request struct {
		// 原始请求结构
		// original request structure
		Raw any

		// 请求头
		Header Header

		// 请求体
		Body io.Reader

		// 请求方法
		Method string

		// 接口定义的虚拟路径
		// The virtual path defined in the interface
		VPath string

		// 实际请求的真实路径
		// Real Path
		RPath string
	}

	// Context 请求上下文
	Context struct {
		// 中间件游标
		// middleware cursor
		index int

		// 配置
		conf *config

		// 缓存
		// session storage
		storage Any

		// 中间件
		// handler chains
		handlers []HandlerFunc

		// 请求
		Request *Request

		// 响应写入器
		Writer ResponseWriter
	}
)

// Close 关闭资源
// close the resource
func Close(resource any) {
	if v, ok := resource.(io.Closer); ok {
		_ = v.Close()
		return
	}
	if v, ok := resource.(Closer); ok {
		v.Close()
	}
}

func NewContext(router *Router, request *Request, writer ResponseWriter) *Context {
	return &Context{
		index:    0,
		conf:     router.conf,
		storage:  Any{},
		handlers: []HandlerFunc{},
		Request:  request,
		Writer:   writer,
	}
}

// Close 关闭请求, 回收Header和Body资源
func (c *Request) Close() {
	Close(c.Body)
	c.Body = nil
}

// Next 执行下一个中间件
func (c *Context) Next() {
	c.index++
	if c.index <= len(c.handlers) {
		c.handlers[c.index-1](c)
	}
}

// Set 设置缓存内容
func (c *Context) Set(key string, val any) {
	c.storage[key] = val
}

// Get 获取缓存内容
func (c *Context) Get(key string) (any, bool) {
	val, ok := c.storage[key]
	return val, ok
}

// WriteJSON 写入JSON
func (c *Context) WriteJSON(code int, v any) error {
	if c.Writer.Protocol() == constant.ProtocolHTTP {
		c.Writer.Header().Set(constant.ContentType, constant.MimeJson)
	}
	c.Writer.Code(code)
	if err := c.conf.jsonCodec.NewEncoder(c.Writer).Encode(v); err != nil {
		return err
	}
	return c.Writer.Flush()
}

// WriteBytes 写入字节流
func (c *Context) WriteBytes(code int, p []byte) error {
	c.Writer.Code(code)
	if _, err := c.Writer.Write(p); err != nil {
		return err
	}
	return c.Writer.Flush()
}

// WriteString 写入字节流
func (c *Context) WriteString(code int, s string) error {
	return c.WriteBytes(code, []byte(s))
}

// BindJSON 绑定请求数据
func (c *Context) BindJSON(v any) error {
	if r, ok := c.Request.Body.(BytesReader); ok {
		return c.conf.jsonCodec.Decode(r.Bytes(), v)
	}
	return c.conf.jsonCodec.NewDecoder(c.Request.Body).Decode(v)
}

// Param 获取路径中的参数
// get the parameters in the path
func (c *Context) Param(key string) string {
	var list1 = internal.Split(c.Request.VPath)
	var list2 = internal.Split(c.Request.RPath)
	var m = len(list1)
	var n = len(list2)
	if m != n {
		return ""
	}
	for i, v := range list1 {
		if isVar(v) && v[1:] == key {
			return list2[i]
		}
	}
	return ""
}

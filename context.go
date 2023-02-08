package uRouter

import (
	"bytes"
	"errors"
	"github.com/lxzan/uRouter/helper"
	"github.com/lxzan/uRouter/internal"
	"io"
	"strings"
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
		Raw() interface{}
	}

	// Request 请求
	Request struct {
		// 原始请求结构
		// original request structure
		Raw interface{}

		// 请求头
		Header Header

		// 请求体
		Body io.Reader

		// 接口中定义的虚拟路径
		// The virtual path defined in the interface
		VPath string
	}

	// Context 请求上下文
	Context struct {
		// 中间件游标
		// middleware cursor
		index int

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
func Close(resource interface{}) {
	if v, ok := resource.(io.Closer); ok {
		_ = v.Close()
		return
	}
	if v, ok := resource.(Closer); ok {
		v.Close()
	}
}

func NewContext(request *Request, writer ResponseWriter) *Context {
	return &Context{
		index:    -1,
		storage:  Any{},
		handlers: []HandlerFunc{},
		Request:  request,
		Writer:   writer,
	}
}

// Next 执行下一个中间件
func (c *Context) Next() {
	c.index++
	if c.index < len(c.handlers) {
		c.handlers[c.index](c)
	}
}

// Set 设置缓存内容
func (c *Context) Set(key string, val interface{}) {
	c.storage[key] = val
}

// Get 获取缓存内容
func (c *Context) Get(key string) (interface{}, bool) {
	val, ok := c.storage[key]
	return val, ok
}

// WriteJSON 写入JSON
func (c *Context) WriteJSON(code int, v interface{}) error {
	if c.Writer.Protocol() == ProtocolHTTP {
		c.Writer.Header().Set(ContentType, MimeJson)
	}
	c.Writer.Code(code)
	if err := JsonCodec().NewEncoder(c.Writer).Encode(v); err != nil {
		return err
	}
	return c.Writer.Flush()
}

// WriteBytes 写入字节流
func (c *Context) WriteBytes(code int, p []byte) error {
	return c.WriteReader(code, bytes.NewBuffer(p))
}

// WriteString 写入字节流
func (c *Context) WriteString(code int, s string) error {
	return c.WriteReader(code, strings.NewReader(s))
}

// WriteReader 写入Reader
func (c *Context) WriteReader(code int, r io.Reader) (err error) {
	c.Writer.Code(code)
	if v, ok := r.(BytesReader); ok {
		err = internal.Write(c.Writer, v.Bytes())
	} else {
		err = internal.Copy(c.Writer, r)
	}
	if err != nil {
		return err
	}
	return c.Writer.Flush()
}

// BindJSON 绑定请求数据
func (c *Context) BindJSON(v interface{}) error {
	defer Close(c.Request.Body)
	if c.Request.Body != nil {
		return defaultJsonCodec.NewDecoder(c.Request.Body).Decode(v)
	}
	return errors.New("request body cannot be nil")
}

// Param 获取路径中的参数
// get the parameters in the path
func (c *Context) Param(key string) string {
	var list1 = helper.Split(c.Request.VPath, defaultSeparator)
	var list2 = helper.Split(c.Request.Header.Get(XPath), defaultSeparator)
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

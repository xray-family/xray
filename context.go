package uRouter

import (
	"bytes"
	"errors"
	"github.com/lxzan/uRouter/internal"
	"io"
	"strings"
)

// TODO buffer pool
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
		Raw    interface{}
		Header Header
		Body   io.Reader
	}

	// Context 请求上下文
	Context struct {
		// 中间件游标
		index int
		// 缓存
		storage Any
		// 中间件
		handlers []HandlerFunc
		// 请求
		Request *Request
		// 响应写入器
		Writer ResponseWriter
	}
)

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
	if err := defaultJsonCodec.NewEncoder(c.Writer).Encode(v); err != nil {
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
	if v, ok := r.(internal.BytesReader); ok {
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
	defer func() {
		_ = internal.Close(c.Request.Body)
	}()
	if c.Request.Body != nil {
		return defaultJsonCodec.NewDecoder(c.Request.Body).Decode(v)
	}
	return errors.New("request body cannot be nil")
}

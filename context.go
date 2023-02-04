package uRouter

import (
	"github.com/lxzan/uRouter/internal"
	"github.com/ugorji/go/codec"
	"io"
	"strings"
)

type (
	// ResponseWriter 响应写入器
	ResponseWriter interface {
		Header() Header

		Write([]byte) (int, error)

		Code(code int)

		// Flush 刷新缓冲区
		Flush() error

		// RawResponseWriter 获取原生响应写入器
		RawResponseWriter() interface{}
	}

	// Request 请求
	Request struct {
		RawRequest interface{}
		Header     Header
		Body       io.Reader
	}

	// Context 请求上下文
	Context struct {
		// 中间件游标
		index int
		// 缓存
		storage A
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
		storage:  A{},
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

// RawRequest 获取原生请求
func (c *Context) RawRequest() interface{} {
	return c.Request.RawRequest
}

// RawResponseWriter 获取原生响应写入器
func (c *Context) RawResponseWriter() interface{} {
	return c.Writer.RawResponseWriter()
}

// WriteJSON 写入JSON
func (c *Context) WriteJSON(code int, v interface{}) error {
	c.Writer.Header().Set(ContentType, MimeJson)
	c.Writer.Code(code)
	if err := codec.NewEncoder(c.Writer, JsonHandle).Encode(v); err != nil {
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
	return c.WriteReader(code, strings.NewReader(s))
}

// WriteReader 写入Reader
func (c *Context) WriteReader(code int, r io.Reader) error {
	c.Writer.Code(code)
	if _, err := io.Copy(c.Writer, r); err != nil {
		return err
	}
	return c.Writer.Flush()
}

// BindJSON 绑定请求数据
func (c *Context) BindJSON(v interface{}) error {
	defer func() {
		_ = internal.Close(c.Request.Body)
	}()

	if br, ok := c.Request.Body.(internal.BytesReader); ok {
		return codec.NewDecoderBytes(br.Bytes(), JsonHandle).Decode(v)
	}
	return codec.NewDecoder(c.Request.Body, JsonHandle).Decode(v)
}

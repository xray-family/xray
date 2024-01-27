package gws

import (
	"bytes"
	"github.com/lxzan/gws"
	"github.com/lxzan/xray"
)

type (
	websocketWrapper interface {
		Writev(opcode gws.Opcode, payloads ...[]byte) error
	}

	responseWriter struct {
		conn        websocketWrapper
		headerCodec *xray.HeaderCodec
		header      xray.Header
		code        gws.Opcode
		buf         *bytes.Buffer
		payloads    [][]byte
	}
)

func (c *responseWriter) Reset(socket websocketWrapper, code gws.Opcode, codec *xray.HeaderCodec) {
	c.conn = socket
	c.code = code
	c.headerCodec = codec
	c.header = codec.Generate()
}

func (c *responseWriter) Protocol() string {
	return xray.ProtocolWebSocket
}

func (c *responseWriter) Raw() interface{} {
	return c.conn
}

func (c *responseWriter) Header() xray.Header {
	return c.header
}

func (c *responseWriter) Code(opcode int) {
	c.code = gws.Opcode(opcode)
}

func (c *responseWriter) RawResponseWriter() interface{} {
	return c.conn
}

func (c *responseWriter) Write(p []byte) (n int, err error) {
	c.payloads = append(c.payloads, p)
	return len(p), nil
}

func (c *responseWriter) Flush() error {
	if err := c.headerCodec.Encode(c.buf, c.header); err != nil {
		return err
	}
	c.payloads[0] = c.buf.Bytes()
	return c.conn.Writev(c.code, c.payloads...)
}

func NewAdapter(router *xray.Router) *Adapter {
	return &Adapter{
		router: router,
		codec:  xray.TextMapHeader,
		pool:   newWriterPool(),
	}
}

type Adapter struct {
	router *xray.Router
	codec  *xray.HeaderCodec
	pool   *writerPool
}

// SetHeaderCodec 设置头部编码方式
func (c *Adapter) SetHeaderCodec(codec *xray.HeaderCodec) *Adapter {
	c.codec = codec
	return c
}

// ServeWebSocket 服务WebSocket
func (c *Adapter) ServeWebSocket(socket *gws.Conn, message *gws.Message) error {
	r := &xray.Request{
		Raw:    message,
		Body:   message,
		Action: "",
	}
	writer := c.pool.Get()
	writer.Reset(socket, message.Opcode, c.codec)
	ctx := xray.NewContext(r, writer)

	header, err := c.codec.Decode(message.Data)
	if err != nil {
		return err
	}

	r.Action = header.Get(xray.XMethod)
	ctx.Request.Header = header
	c.router.EmitEvent(r.Action, header.Get(xray.XPath), ctx)
	return nil
}

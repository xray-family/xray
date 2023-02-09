package gws

import (
	"bytes"
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
	"github.com/lxzan/uRouter/constant"
	"sync"
)

type (
	websocketWrapper interface {
		WriteMessage(opcode gws.Opcode, payload []byte)
	}

	responseWriter struct {
		once        sync.Once
		conn        websocketWrapper
		headerCodec *uRouter.HeaderCodec
		header      uRouter.Header
		code        gws.Opcode
		buf         *bytes.Buffer
	}
)

func newResponseWriter(socket websocketWrapper, codec *uRouter.HeaderCodec) *responseWriter {
	return &responseWriter{
		once:        sync.Once{},
		code:        gws.OpcodeText,
		conn:        socket,
		headerCodec: codec,
		header:      codec.Generate(),
		buf:         uRouter.BufferPool().Get(),
	}
}

func (c *responseWriter) Protocol() string {
	return uRouter.ProtocolWebSocket
}

func (c *responseWriter) Raw() interface{} {
	return c.conn
}

func (c *responseWriter) Header() uRouter.Header {
	return c.header
}

func (c *responseWriter) Code(opcode int) {
	c.code = gws.Opcode(opcode)
}

func (c *responseWriter) RawResponseWriter() interface{} {
	return c.conn
}

func (c *responseWriter) Write(p []byte) (n int, err error) {
	c.once.Do(func() {
		err = c.headerCodec.Encode(c.buf, c.header)
		n = c.buf.Len()
	})
	if err != nil {
		return
	}
	return c.buf.Write(p)
}

func (c *responseWriter) Flush() error {
	c.conn.WriteMessage(c.code, c.buf.Bytes())
	uRouter.BufferPool().Put(c.buf)
	return nil
}

func NewAdapter() *Adapter {
	return &Adapter{
		Router: uRouter.New(),
		codec:  uRouter.TextHeader,
	}
}

type Adapter struct {
	*uRouter.Router
	codec *uRouter.HeaderCodec
}

// SetHeaderCodec 设置头部编码方式
func (c *Adapter) SetHeaderCodec(codec *uRouter.HeaderCodec) *Adapter {
	c.codec = codec
	return c
}

// ServeWebSocket 服务WebSocket
func (c *Adapter) ServeWebSocket(socket *gws.Conn, message *gws.Message) error {
	r := &uRouter.Request{
		Raw:    message,
		Body:   message,
		Action: "",
	}
	ctx := uRouter.NewContext(r, newResponseWriter(socket, c.codec))

	header, err := c.codec.Decode(message.Data)
	if err != nil {
		return err
	}

	ctx.Request.Header = header
	c.Router.EmitAction(r.Action, header.Get(constant.XPath), ctx)
	return nil
}

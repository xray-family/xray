package gws

import (
	"bytes"
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
)

type (
	websocket interface {
		WriteMessage(opcode gws.Opcode, payload []byte)
	}

	responseWriter struct {
		conn        websocket
		headerCodec *uRouter.HeaderCodec
		header      uRouter.Header
		code        gws.Opcode
		payload     []byte
	}
)

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

func (c *responseWriter) Write(p []byte) (int, error) {
	c.payload = p
	return len(p), nil
}

func (c *responseWriter) Flush() error {
	if c.code == 0 {
		c.code = gws.OpcodeText
	}
	var buf = bytes.NewBuffer(nil)
	if err := c.headerCodec.Encode(buf, c.header); err != nil {
		return err
	}
	buf.Write(c.payload)
	c.conn.WriteMessage(c.code, buf.Bytes())
	return nil
}

func NewAdapter(r *uRouter.Router, codec *uRouter.HeaderCodec) *Adapter {
	return &Adapter{
		router: r,
		codec:  codec,
	}
}

type Adapter struct {
	router *uRouter.Router
	codec  *uRouter.HeaderCodec
}

func (c *Adapter) ServeWebSocket(socket *gws.Conn, p []byte) error {
	message := bytes.NewBuffer(p)
	ctx := &uRouter.Context{
		Request: &uRouter.Request{
			Body: message,
		},
		Writer: &responseWriter{
			conn:        socket,
			headerCodec: c.codec,
			header:      c.codec.Generate(),
		},
	}

	header, err := c.codec.Decode(message)
	if err != nil {
		return err
	}

	ctx.Request.Header = header
	c.router.Emit(ctx)
	return nil
}

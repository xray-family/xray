package gws

import (
	"bytes"
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
)

type Writer struct {
	conn        *gws.Conn
	headerCodec *uRouter.HeaderCodec
	header      uRouter.Header
	code        gws.Opcode
	payload     []byte
}

func (c *Writer) Header() uRouter.Header {
	return c.header
}

func (c *Writer) Code(opcode int) {
	c.code = gws.Opcode(opcode)
}

func (c *Writer) RawResponseWriter() interface{} {
	return c.code
}

func (c *Writer) Write(p []byte) (int, error) {
	c.payload = p
	return len(p), nil
}

func (c *Writer) Flush() error {
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
		Writer: &Writer{
			conn:        socket,
			headerCodec: c.codec,
			header:      c.codec.Generate(),
		},
	}

	header, err := c.router.HeaderCodec.Decode(message)
	if err != nil {
		return err
	}

	ctx.Request.Header = header
	c.router.Emit(ctx)
	return nil
}

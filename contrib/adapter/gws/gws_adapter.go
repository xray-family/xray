package gws

import (
	"bytes"
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
	"github.com/lxzan/uRouter/internal"
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
		headerBuf   *bytes.Buffer
		bodyBuf     *bytes.Buffer
	}
)

func newResponseWriter(socket websocket, codec *uRouter.HeaderCodec) *responseWriter {
	return &responseWriter{
		code:        gws.OpcodeText,
		conn:        socket,
		headerCodec: codec,
		header:      codec.Generate(),
		headerBuf:   bytes.NewBuffer(make([]byte, 0, 1024)),
		bodyBuf:     bytes.NewBuffer(make([]byte, 0, 1024)),
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

func (c *responseWriter) Write(p []byte) (int, error) {
	return c.bodyBuf.Write(p)
}

func (c *responseWriter) Flush() error {
	if c.code == 0 {
		c.code = gws.OpcodeText
	}

	if err := c.headerCodec.Encode(c.headerBuf, c.header); err != nil {
		return err
	}
	c.headerBuf.Write(c.bodyBuf.Bytes())
	c.conn.WriteMessage(c.code, c.headerBuf.Bytes())
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

func (c *Adapter) ServeWebSocket(socket *gws.Conn, message internal.BytesReader) error {
	ctx := &uRouter.Context{
		Request: &uRouter.Request{
			Raw:  message,
			Body: message,
		},
		Writer: newResponseWriter(socket, c.codec),
	}

	header, err := c.codec.Decode(message)
	if err != nil {
		return err
	}

	ctx.Request.Header = header
	c.router.Emit(ctx)
	return nil
}

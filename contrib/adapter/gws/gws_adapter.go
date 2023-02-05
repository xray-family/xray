package gws

import (
	"bytes"
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
	"github.com/lxzan/uRouter/internal"
	"sync"
)

type (
	websocket interface {
		WriteMessage(opcode gws.Opcode, payload []byte)
	}

	responseWriter struct {
		once        sync.Once
		conn        websocket
		headerCodec *uRouter.HeaderCodec
		header      uRouter.Header
		code        gws.Opcode
		buf         *bytes.Buffer
	}
)

func newResponseWriter(socket websocket, codec *uRouter.HeaderCodec) *responseWriter {
	return &responseWriter{
		once:        sync.Once{},
		code:        gws.OpcodeText,
		conn:        socket,
		headerCodec: codec,
		header:      codec.Generate(),
		buf:         uRouter.DefaultBufferPool.Get(),
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
	uRouter.DefaultBufferPool.Put(c.buf)
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

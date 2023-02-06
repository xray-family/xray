package gorilla

import (
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/lxzan/uRouter"
	"sync"
)

type (
	Message struct {
		Opcode int
		Data   *bytes.Buffer
	}

	websocketWrapper interface {
		WriteMessage(opcode int, payload []byte) error
	}

	responseWriter struct {
		once        sync.Once
		conn        websocketWrapper
		headerCodec *uRouter.HeaderCodec
		header      uRouter.Header
		code        int
		buf         *bytes.Buffer
	}
)

func (c *Message) Bytes() []byte {
	return c.Data.Bytes()
}

func (c *Message) Read(p []byte) (n int, err error) {
	return c.Data.Read(p)
}

func newResponseWriter(socket websocketWrapper, codec *uRouter.HeaderCodec) *responseWriter {
	return &responseWriter{
		once:        sync.Once{},
		code:        websocket.TextMessage,
		conn:        socket,
		headerCodec: codec,
		header:      codec.Generate(),
		buf:         uRouter.DefaultBufferPool().Get(),
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
	c.code = int(opcode)
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
	uRouter.DefaultBufferPool().Put(c.buf)
	return nil
}

func NewAdapter(r *uRouter.Router) *Adapter {
	return &Adapter{
		router: r,
		codec:  uRouter.TextHeader,
	}
}

type Adapter struct {
	router *uRouter.Router
	codec  *uRouter.HeaderCodec
}

// SetHeaderCodec 设置头部编码方式
func (c *Adapter) SetHeaderCodec(codec *uRouter.HeaderCodec) *Adapter {
	c.codec = codec
	return c
}

func (c *Adapter) ServeWebSocket(socket *websocket.Conn, opcode int, p []byte) error {
	var message = &Message{
		Opcode: opcode,
		Data:   bytes.NewBuffer(p),
	}
	ctx := uRouter.NewContext(
		&uRouter.Request{Raw: message, Body: message},
		newResponseWriter(socket, c.codec),
	)

	header, err := c.codec.Decode(message)
	if err != nil {
		return err
	}

	ctx.Request.Header = header
	c.router.Emit(ctx)
	return nil

}

package gws

import (
	"bytes"
	"fmt"
	"github.com/lxzan/gws"
	"github.com/xray-family/xray"
	"github.com/xray-family/xray/codec"
	"github.com/xray-family/xray/internal"
	"io"
	"strconv"
)

type (
	websocketWrapper interface {
		Writev(opcode gws.Opcode, payloads ...[]byte) error
	}

	responseWriter struct {
		conn     websocketWrapper
		opcode   gws.Opcode
		codec    codec.Codec
		header   xray.Header
		buf      *bytes.Buffer
		payloads [][]byte
	}
)

func (c *responseWriter) Protocol() string {
	return xray.ProtocolWebSocket
}

func (c *responseWriter) Raw() any {
	return c.conn
}

func (c *responseWriter) Header() xray.Header {
	return c.header
}

func (c *responseWriter) Code(code int) {
	c.header.Set(xray.XStatus, strconv.Itoa(code))
}

func (c *responseWriter) RawResponseWriter() any {
	return c.conn
}

func (c *responseWriter) Write(p []byte) (n int, err error) {
	c.payloads = append(c.payloads, internal.Clone(p))
	return len(p), nil
}

func (c *responseWriter) Flush() error {
	if err := marshalHeader(c.codec, c.buf, c.header); err != nil {
		return err
	}
	c.payloads[0] = c.buf.Bytes()
	return c.conn.Writev(c.opcode, c.payloads...)
}

func NewAdapter(router *xray.Router) *Adapter {
	return &Adapter{
		router: router,
		tpl:    &xray.SliceHeader{},
		pool:   newWriterPool(),
	}
}

type Adapter struct {
	router *xray.Router
	pool   *writerPool
	tpl    xray.Header
}

// SetHeaderTpl 设置头部编码方式
func (c *Adapter) SetHeaderTpl(tpl xray.Header) *Adapter {
	c.tpl = tpl
	return c
}

// ServeWebSocket 服务WebSocket
func (c *Adapter) ServeWebSocket(socket *gws.Conn, message *gws.Message) error {
	r := &xray.Request{Raw: message, Body: message}
	w := c.pool.Get()
	w.conn = socket
	w.opcode = message.Opcode
	w.codec = c.router.JsonCodec()
	w.header = c.tpl.New()
	w.payloads = append(w.payloads, nil)
	ctx := xray.NewContext(c.router, r, w)

	header := c.tpl.New()
	if err := unmarshalHeader(c.router.JsonCodec(), message.Data, header); err != nil {
		return err
	}

	r.Method = header.Get(xray.XMethod)
	ctx.Request.Header = header
	c.router.EmitEvent(r.Method, header.Get(xray.XPath), ctx)
	c.pool.Put(w)
	return nil
}

func marshalHeader(jsonCodec codec.Codec, w *bytes.Buffer, v xray.Header) error {
	if v.Len() == 0 {
		w.WriteString("0002[]")
		return nil
	}
	w.WriteString("0000")
	if err := jsonCodec.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	copy(w.Bytes()[:4], fmt.Sprintf("%04d", w.Len()-4))
	return nil
}

func unmarshalHeader(jsonCodec codec.Codec, r *bytes.Buffer, v xray.Header) error {
	if r.Len() < 4 {
		return io.ErrShortBuffer
	}
	length, err := strconv.Atoi(string(r.Next(4)))
	if err != nil {
		return err
	}
	p := r.Next(length)
	if len(p) < length {
		return io.ErrShortBuffer
	}
	return jsonCodec.Decode(p, v)
}

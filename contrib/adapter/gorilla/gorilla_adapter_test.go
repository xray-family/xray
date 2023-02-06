package gorilla

import (
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/lxzan/uRouter"
	"github.com/stretchr/testify/assert"
	"testing"
)

type connMocker struct {
	opcode int
	buf    *bytes.Buffer
}

func (c *connMocker) WriteMessage(opcode int, payload []byte) error {
	c.opcode = opcode
	c.buf.Write(payload)
	return nil
}

type messageMocker struct {
	b *bytes.Buffer
}

func (c *messageMocker) Read(p []byte) (n int, err error) {
	return c.b.Read(p)
}

func (c *messageMocker) Bytes() []byte {
	return c.b.Bytes()
}

func TestNewAdapter(t *testing.T) {
	var as = assert.New(t)

	t.Run("normal", func(t *testing.T) {
		const requestPayload = "hello"
		const responsePayload = "world"
		var sum = 0
		var router = uRouter.New()
		var adapter = NewAdapter(router).SetHeaderCodec(uRouter.TextHeader)

		router.On("testEncode", func(ctx *uRouter.Context) {
			ctx.Writer = newResponseWriter(&connMocker{buf: bytes.NewBufferString("")}, uRouter.TextHeader)

			sum++
			ctx.Writer.Header().Set(uRouter.ContentType, uRouter.MimeStream)
			ctx.Writer.Header().Set(uRouter.XPath, "/testDecode")
			ctx.Writer.Code(websocket.TextMessage)
			_, _ = ctx.Writer.Write([]byte(responsePayload))
			ctx.Writer.Raw()
			if err := ctx.Writer.Flush(); err != nil {
				as.NoError(err)
				return
			}

			as.Equal(2, ctx.Request.Header.Len())
			as.Equal(requestPayload, string(ctx.Request.Body.(*Message).Bytes()))

			var writer = ctx.Writer.Raw().(*connMocker)
			if err := adapter.ServeWebSocket(nil, websocket.TextMessage, writer.buf.Bytes()); err != nil {
				as.NoError(err)
				return
			}
		})

		router.On("testDecode", func(ctx *uRouter.Context) {
			sum++
			as.Equal(2, ctx.Request.Header.Len())
			as.Equal(responsePayload, string(ctx.Request.Body.(*Message).Bytes()))
		})

		var b = &messageMocker{b: bytes.NewBufferString("")}
		var header = uRouter.MapHeader{
			uRouter.ContentType: uRouter.MimeJson,
			uRouter.XPath:       "/testEncode",
		}
		if err := adapter.codec.Encode(b.b, header); err != nil {
			as.NoError(err)
			return
		}
		b.b.WriteString(requestPayload)
		if err := adapter.ServeWebSocket(nil, websocket.TextMessage, b.Bytes()); err != nil {
			as.NoError(err)
			return
		}
		as.Equal(sum, 2)
	})
}

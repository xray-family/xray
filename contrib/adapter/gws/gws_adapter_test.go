package gws

import (
	"bytes"
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
	"github.com/lxzan/uRouter/constant"
	"github.com/stretchr/testify/assert"
	"testing"
)

type connMocker struct {
	opcode gws.Opcode
	buf    *bytes.Buffer
}

func (c *connMocker) WriteMessage(opcode gws.Opcode, payload []byte) {
	c.opcode = opcode
	c.buf.Write(payload)
}

func TestNewAdapter(t *testing.T) {
	var as = assert.New(t)

	t.Run("normal", func(t *testing.T) {
		const requestPayload = "hello"
		const responsePayload = "world"
		var sum = 0
		var router = uRouter.New()
		var adapter = NewAdapter(router).SetHeaderCodec(uRouter.TextMapHeader)

		router.On("testEncode", func(ctx *uRouter.Context) {
			ctx.Writer = newResponseWriter(&connMocker{buf: bytes.NewBufferString("")}, uRouter.TextMapHeader)

			sum++
			ctx.Writer.Header().Set(constant.ContentType, constant.MimeStream)
			ctx.Writer.Header().Set(constant.XPath, "/testDecode")
			ctx.Writer.Code(int(gws.OpcodeText))
			_, _ = ctx.Writer.Write([]byte(responsePayload))
			ctx.Writer.Raw()
			if err := ctx.Writer.Flush(); err != nil {
				as.NoError(err)
				return
			}

			as.Equal(2, ctx.Request.Header.Len())
			as.Equal(requestPayload, ctx.Request.Body.(*gws.Message).Data.String())

			var writer = ctx.Writer.Raw().(*connMocker)
			if err := adapter.ServeWebSocket(nil, &gws.Message{Data: writer.buf}); err != nil {
				as.NoError(err)
				return
			}
		})

		router.On("testDecode", func(ctx *uRouter.Context) {
			sum++
			as.Equal(2, ctx.Request.Header.Len())
			as.Equal(responsePayload, ctx.Request.Body.(*gws.Message).Data.String())
		})

		var b = &gws.Message{
			Opcode: gws.OpcodeText,
			Data:   bytes.NewBufferString(""),
		}
		var header = uRouter.NewMapHeader()
		header.Set(constant.ContentType, constant.MimeJson)
		header.Set(constant.XPath, "/testEncode")
		if err := adapter.codec.Encode(b.Data, header); err != nil {
			as.NoError(err)
			return
		}
		b.Data.WriteString(requestPayload)
		if err := adapter.ServeWebSocket(nil, b); err != nil {
			as.NoError(err)
			return
		}
		as.Equal(sum, 2)
	})
}

func TestOthers(t *testing.T) {
	var w = newResponseWriter(&gws.Conn{}, uRouter.TextMapHeader)
	assert.Equal(t, uRouter.ProtocolWebSocket, w.Protocol())

	w.RawResponseWriter()
}

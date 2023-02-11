package gws

import (
	"bytes"
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
	"github.com/lxzan/uRouter/constant"
	"github.com/lxzan/uRouter/internal"
	"testing"
)

func BenchmarkResponseWriter_Write1024(b *testing.B) {
	ctx := uRouter.NewContext(
		&uRouter.Request{Header: uRouter.TextMapHeader.Generate()},
		newResponseWriter(&connMocker{
			buf: bytes.NewBuffer(make([]byte, 0, constant.BufferLeveL16)),
		}, uRouter.TextMapHeader),
	)

	var v = struct {
		Message string `json:"message"`
	}{
		Message: string(internal.AlphabetNumeric.Generate(1024)),
	}

	for i := 0; i < b.N; i++ {
		ctx.Request.Header.Set(constant.XRealIP, "127.0.0.1")
		ctx.Request.Header.Set(uRouter.UPath, "/test")
		_ = ctx.WriteJSON(int(gws.OpcodeText), &v)

		ctx.Request.Header = uRouter.TextMapHeader.Generate()
		ctx.Writer.Raw().(*connMocker).buf.Reset()
	}
}

func BenchmarkResponseWriter_Write512(b *testing.B) {
	ctx := uRouter.NewContext(
		&uRouter.Request{Header: uRouter.TextMapHeader.Generate(), Raw: &gws.Message{}},
		newResponseWriter(&connMocker{
			buf: bytes.NewBuffer(make([]byte, 0, constant.BufferLeveL16)),
		}, uRouter.TextMapHeader),
	)

	var v = struct {
		Message string `json:"message"`
	}{
		Message: string(internal.AlphabetNumeric.Generate(16)),
	}

	for i := 0; i < b.N; i++ {
		ctx.Request.Header = uRouter.HeaderPool().Get(constant.MapHeaderNumber)
		ctx.Request.Header.Set(constant.XRealIP, "127.0.0.1")
		ctx.Request.Header.Set(uRouter.UPath, "/test")
		_ = ctx.WriteJSON(int(gws.OpcodeText), &v)

		uRouter.HeaderPool().Put(ctx.Request.Header)
		ctx.Writer.Raw().(*connMocker).buf.Reset()
	}
}

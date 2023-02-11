package http

import (
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
	"github.com/lxzan/uRouter/internal"
	"net/http"
	"testing"
)

func BenchmarkHTTPReadWrite16(b *testing.B) {
	ctx := uRouter.NewContext(
		&uRouter.Request{Header: uRouter.NewHttpHeader(http.Header{})},
		&responseWriter{ResponseWriter: newWriterMocker()},
	)

	var v = struct {
		Message string `json:"message"`
	}{
		Message: string(internal.AlphabetNumeric.Generate(16)),
	}

	for i := 0; i < b.N; i++ {
		ctx.Request.Header = uRouter.HeaderPool().Get(internal.MapHeaderNumber)
		ctx.Request.Header.Set(internal.XRealIP, "127.0.0.1")
		ctx.Request.Header.Set(uRouter.UPath, "/test")
		_ = ctx.WriteJSON(int(gws.OpcodeText), &v)

		ctx.Writer.Raw().(*writerMocker).buf.Reset()
	}
}

func BenchmarkHTTPReadWrite512(b *testing.B) {
	r := uRouter.New()
	path := "/test"
	r.OnEvent(http.MethodGet, path, func(ctx *uRouter.Context) {
	})
	r.Start()

	ctx := uRouter.NewContext(
		&uRouter.Request{Header: uRouter.NewHttpHeader(http.Header{})},
		&responseWriter{ResponseWriter: newWriterMocker()},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.EmitEvent(http.MethodGet, path, ctx)
	}
}

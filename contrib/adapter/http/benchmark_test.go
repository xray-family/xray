package http

import (
	_ "embed"
	"github.com/lxzan/uRouter"
	"github.com/lxzan/uRouter/codec"
	"github.com/lxzan/uRouter/internal"
	"net/http"
	"net/url"
	"testing"
)

func BenchmarkStaticRoute(b *testing.B) {
	r := uRouter.New()
	path := "/test"
	r.OnEvent(http.MethodGet, path, func(ctx *uRouter.Context) {
	})
	r.StartSilently()

	ctx := uRouter.NewContext(
		&uRouter.Request{Header: uRouter.NewHttpHeader(http.Header{})},
		&responseWriter{ResponseWriter: newWriterMocker()},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.EmitEvent(http.MethodGet, path, ctx)
	}
}

func BenchmarkDynamicRoute(b *testing.B) {
	r := uRouter.New()
	path := "/api/v1/user/:id"
	r.OnEvent(http.MethodGet, path, func(ctx *uRouter.Context) {
	})
	r.StartSilently()

	ctx := uRouter.NewContext(
		&uRouter.Request{Header: uRouter.NewHttpHeader(http.Header{})},
		&responseWriter{ResponseWriter: newWriterMocker()},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.EmitEvent(http.MethodGet, "/api/v1/user/1", ctx)
	}
}

//go:embed github-api.json
var githubAPI []byte

func BenchmarkGithubRoute(b *testing.B) {
	r := uRouter.New()
	var rows = uRouter.Form{}
	codec.StdJsonCodec.Decode(githubAPI, &rows)

	var paths []string
	for _, s := range rows {
		u, e := url.Parse(s)
		if e == nil {
			paths = append(paths, u.Path)
			r.OnEvent(http.MethodGet, u.Path, func(ctx *uRouter.Context) {})
		}
	}
	r.StartSilently()

	ctx := uRouter.NewContext(
		&uRouter.Request{Header: uRouter.NewHttpHeader(http.Header{})},
		&responseWriter{ResponseWriter: newWriterMocker()},
	)

	b.ResetTimer()
	count := len(paths)
	for i := 0; i < b.N; i++ {
		path := paths[i%count]
		r.EmitEvent(http.MethodGet, path, ctx)
	}
}

func BenchmarkWrite(b *testing.B) {
	var message = struct {
		Message string `json:"message"`
	}{
		Message: string(internal.AlphabetNumeric.Generate(512)),
	}

	r := uRouter.New()
	path := "/test"
	r.OnEvent(http.MethodGet, path, func(ctx *uRouter.Context) {
		ctx.Writer.Raw().(*writerMocker).buf.Reset()
		_ = ctx.WriteJSON(http.StatusOK, message)
	})
	r.StartSilently()

	ctx := uRouter.NewContext(
		&uRouter.Request{Header: uRouter.NewHttpHeader(http.Header{})},
		&responseWriter{ResponseWriter: newWriterMocker()},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.EmitEvent(http.MethodGet, path, ctx)
	}
}

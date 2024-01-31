package http

import (
	"bytes"
	_ "embed"
	"github.com/xray-family/xray"
	"net/http"
	"testing"
)

func BenchmarkAdapter_ServeHTTP(b *testing.B) {
	var router = xray.New()
	router.GET("/", func(ctx *xray.Context) {})
	adapter := NewAdapter(router)
	w := &responseWriter{ResponseWriter: newWriterMocker()}
	r, _ := http.NewRequest(http.MethodGet, "http://example.com/", bytes.NewBufferString(""))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		adapter.ServeHTTP(w.ResponseWriter, r)
	}
}

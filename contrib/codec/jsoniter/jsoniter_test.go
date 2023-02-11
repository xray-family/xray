package jsoniter

import (
	"github.com/lxzan/uRouter"
	"github.com/lxzan/uRouter/codec"
	"github.com/lxzan/uRouter/internal"
	"testing"
)

func BenchmarkJsoniterCodec(b *testing.B) {
	var s = string(internal.AlphabetNumeric.Generate(internal.BufferLeveL1))
	var v = struct {
		Message string
	}{Message: s}

	for i := 0; i < b.N; i++ {
		_, _ = JsoniterCodec.Encode(v)
	}
}

func BenchmarkStdJsonCodec(b *testing.B) {
	var s = string(internal.AlphabetNumeric.Generate(internal.BufferLeveL1))
	var v = struct {
		Message string
	}{Message: s}

	for i := 0; i < b.N; i++ {
		_, _ = codec.StdJsonCodec.Encode(v)
	}
}

func BenchmarkPoolJsonCodec(b *testing.B) {
	var s = string(internal.AlphabetNumeric.Generate(internal.BufferLeveL1 - 64))
	var v = struct {
		Message string
	}{Message: s}

	for i := 0; i < b.N; i++ {
		w := uRouter.BufferPool().Get(internal.BufferLeveL1)
		_ = codec.StdJsonCodec.NewEncoder(w).Encode(v)
		uRouter.BufferPool().Put(w)
	}
}

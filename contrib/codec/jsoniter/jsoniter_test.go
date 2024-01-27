package jsoniter

import (
	"github.com/lxzan/xray"
	"github.com/lxzan/xray/codec"
	"github.com/lxzan/xray/constant"
	"github.com/lxzan/xray/internal"
	"testing"
)

func BenchmarkJsoniterCodec(b *testing.B) {
	var s = string(internal.AlphabetNumeric.Generate(constant.BufferLeveL1))
	var v = struct {
		Message string
	}{Message: s}

	for i := 0; i < b.N; i++ {
		_, _ = JsoniterCodec.Encode(v)
	}
}

func BenchmarkStdJsonCodec(b *testing.B) {
	var s = string(internal.AlphabetNumeric.Generate(constant.BufferLeveL1))
	var v = struct {
		Message string
	}{Message: s}

	for i := 0; i < b.N; i++ {
		_, _ = codec.StdJsonCodec.Encode(v)
	}
}

func BenchmarkPoolJsonCodec(b *testing.B) {
	var s = string(internal.AlphabetNumeric.Generate(constant.BufferLeveL1 - 64))
	var v = struct {
		Message string
	}{Message: s}

	for i := 0; i < b.N; i++ {
		w := xray.BufferPool().Get(constant.BufferLeveL1)
		_ = codec.StdJsonCodec.NewEncoder(w).Encode(v)
		xray.BufferPool().Put(w)
	}
}

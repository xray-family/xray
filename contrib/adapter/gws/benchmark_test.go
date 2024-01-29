package gws

//
//import (
//	"bytes"
//	"github.com/lxzan/gws"
//	"github.com/lxzan/xray"
//	"github.com/lxzan/xray/constant"
//	"github.com/lxzan/xray/internal"
//	"testing"
//)
//
//func BenchmarkAdapter_ServeHTTP(b *testing.B) {
//	var router = xray.New()
//	router.OnGET("/", func(ctx *xray.Context) {})
//	adapter := NewAdapter(router)
//
//	socket := &gws.Conn{}
//	msg := &gws.Message{Data: bytes.NewBufferString("")}
//	header := xray.MapHeaderTemplate.Generate()
//	header.Set(xray.XPath, "/")
//	header.Set(xray.XMethod, "")
//	_ = xray.TextMapHeader.Encode(msg.Data, header)
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_ = adapter.ServeWebSocket(socket, msg)
//	}
//}
//
//func BenchmarkResponseWriter_Write1024(b *testing.B) {
//	ctx := xray.NewContext(
//		&xray.Request{Header: xray.TextMapHeader.Generate()},
//		newResponseWriter(&connMocker{
//			buf: bytes.NewBuffer(make([]byte, 0, constant.BufferLeveL16)),
//		}, xray.TextMapHeader),
//	)
//
//	var v = struct {
//		Message string `json:"message"`
//	}{
//		Message: string(internal.AlphabetNumeric.Generate(1024)),
//	}
//
//	for i := 0; i < b.N; i++ {
//		ctx.Request.Header.Set(XRealIP, "127.0.0.1")
//		ctx.Request.Header.Set(xray.XPath, "/test")
//		_ = ctx.WriteJSON(int(gws.OpcodeText), &v)
//
//		ctx.Request.Header = xray.TextMapHeader.Generate()
//		ctx.Writer.Raw().(*connMocker).buf.Reset()
//	}
//}
//
//func BenchmarkResponseWriter_Write512(b *testing.B) {
//	ctx := xray.NewContext(
//		&xray.Request{Header: xray.TextMapHeader.Generate(), Raw: &gws.Message{}},
//		newResponseWriter(&connMocker{
//			buf: bytes.NewBuffer(make([]byte, 0, constant.BufferLeveL16)),
//		}, xray.TextMapHeader),
//	)
//
//	var v = struct {
//		Message string `json:"message"`
//	}{
//		Message: string(internal.AlphabetNumeric.Generate(16)),
//	}
//
//	for i := 0; i < b.N; i++ {
//		ctx.Request.Header = xray.HeaderPool().Get(constant.MapHeaderNumber)
//		ctx.Request.Header.Set(XRealIP, "127.0.0.1")
//		ctx.Request.Header.Set(xray.XPath, "/test")
//		_ = ctx.WriteJSON(int(gws.OpcodeText), &v)
//
//		xray.HeaderPool().Put(constant.MapHeaderNumber, ctx.Request.Header)
//		ctx.Writer.Raw().(*connMocker).buf.Reset()
//	}
//}

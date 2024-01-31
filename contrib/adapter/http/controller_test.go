package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/xray-family/xray"
	"net/http"
	"testing"
)

func TestRootController_ClientIP(t *testing.T) {
	ctl := new(RootController)

	t.Run("", func(t *testing.T) {
		router := xray.New()
		ctx := xray.NewContext(router, &xray.Request{
			Header: xray.HttpHeader{Header: map[string][]string{}},
		}, nil)
		ctx.Request.Header.Set(xray.XRealIP, "127.0.0.2")
		assert.Equal(t, "127.0.0.2", ctl.ClientIP(ctx))
	})

	t.Run("", func(t *testing.T) {
		router := xray.New()
		ctx := xray.NewContext(router, &xray.Request{
			Raw:    &http.Request{RemoteAddr: "127.0.0.1:3000"},
			Header: xray.HttpHeader{Header: map[string][]string{}},
		}, nil)
		assert.Equal(t, "127.0.0.1", ctl.ClientIP(ctx))
	})

	t.Run("", func(t *testing.T) {
		router := xray.New()
		ctx := xray.NewContext(router, &xray.Request{
			Raw:    &http.Request{RemoteAddr: "[::]:3000"},
			Header: xray.HttpHeader{Header: map[string][]string{}},
		}, nil)
		assert.Equal(t, "127.0.0.1", ctl.ClientIP(ctx))
	})

	t.Run("", func(t *testing.T) {
		router := xray.New()
		ctx := xray.NewContext(router, &xray.Request{
			Raw:    &http.Request{RemoteAddr: "[::]:3000"},
			Header: xray.HttpHeader{Header: map[string][]string{}},
		}, &responseWriter{newWriterMocker()})
		writer := ctl.GetResponseWriter(ctx)
		assert.NotNil(t, writer)
	})
}

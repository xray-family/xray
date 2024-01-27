package http

import (
	"github.com/lxzan/xray"
	"github.com/lxzan/xray/constant"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRootController_ClientIP(t *testing.T) {
	ctl := new(RootController)

	t.Run("", func(t *testing.T) {
		ctx := xray.NewContext(&xray.Request{
			Header: xray.HttpHeader{Header: map[string][]string{}},
		}, nil)
		ctx.Request.Header.Set(constant.XRealIP, "127.0.0.2")
		assert.Equal(t, "127.0.0.2", ctl.ClientIP(ctx))
	})

	t.Run("", func(t *testing.T) {
		ctx := xray.NewContext(&xray.Request{
			Raw:    &http.Request{RemoteAddr: "127.0.0.1:3000"},
			Header: xray.HttpHeader{Header: map[string][]string{}},
		}, nil)
		assert.Equal(t, "127.0.0.1", ctl.ClientIP(ctx))
	})

	t.Run("", func(t *testing.T) {
		ctx := xray.NewContext(&xray.Request{
			Raw:    &http.Request{RemoteAddr: "[::]:3000"},
			Header: xray.HttpHeader{Header: map[string][]string{}},
		}, nil)
		assert.Equal(t, "127.0.0.1", ctl.ClientIP(ctx))
	})
}

package http

import (
	"github.com/lxzan/uRouter"
	"github.com/lxzan/uRouter/internal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRootController_ClientIP(t *testing.T) {
	ctl := new(RootController)

	t.Run("", func(t *testing.T) {
		ctx := uRouter.NewContext(&uRouter.Request{
			Header: uRouter.HttpHeader{Header: map[string][]string{}},
		}, nil)
		ctx.Request.Header.Set(internal.XRealIP, "127.0.0.2")
		assert.Equal(t, "127.0.0.2", ctl.ClientIP(ctx))
	})

	t.Run("", func(t *testing.T) {
		ctx := uRouter.NewContext(&uRouter.Request{
			Raw:    &http.Request{RemoteAddr: "127.0.0.1:3000"},
			Header: uRouter.HttpHeader{Header: map[string][]string{}},
		}, nil)
		assert.Equal(t, "127.0.0.1", ctl.ClientIP(ctx))
	})

	t.Run("", func(t *testing.T) {
		ctx := uRouter.NewContext(&uRouter.Request{
			Raw:    &http.Request{RemoteAddr: "[::]:3000"},
			Header: uRouter.HttpHeader{Header: map[string][]string{}},
		}, nil)
		assert.Equal(t, "127.0.0.1", ctl.ClientIP(ctx))
	})
}

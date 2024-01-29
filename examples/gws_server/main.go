package main

import (
	"github.com/lxzan/gws"
	"github.com/lxzan/xray"
	gwsAdapter "github.com/lxzan/xray/contrib/adapter/gws"
	httpAdapter "github.com/lxzan/xray/contrib/adapter/http"
	"github.com/lxzan/xray/contrib/codec/jsoniter"
	"github.com/lxzan/xray/contrib/log/zerolog"
	"github.com/lxzan/xray/log"
	"net/http"
)

func main() {
	router := xray.New(
		xray.WithJsonCodec(jsoniter.Codec),
		xray.WithLogger(zerolog.Logger),
	)
	router.Use(xray.Recovery())

	upgrader := gws.NewUpgrader(&WebSocketHandler{logger: router.Logger(), adapter: gwsAdapter.NewAdapter(router)}, nil)

	router.OnGET("/connect", func(ctx *xray.Context) {
		socket, err := upgrader.Upgrade(ctx.Writer.Raw().(http.ResponseWriter), ctx.Request.Raw.(*http.Request))
		if err != nil {
			router.Logger().Error(err.Error())
			return
		}
		go socket.ReadLoop()
	})

	router.On("/greet", func(ctx *xray.Context) {
		ctx.Writer.Header().Set("content-type", "plain/text")
		_ = ctx.WriteString(int(gws.OpcodeText), "hello!")
	})

	if err := http.ListenAndServe(":3000", httpAdapter.NewAdapter(router)); err != nil {
		router.Logger().Panic(err.Error())
	}
}

type WebSocketHandler struct {
	gws.BuiltinEventHandler
	logger  log.Logger
	adapter *gwsAdapter.Adapter
}

func (c *WebSocketHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	if err := c.adapter.ServeWebSocket(socket, message); err != nil {
		c.logger.Error(err.Error())
	}
}

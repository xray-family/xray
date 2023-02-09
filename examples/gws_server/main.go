package main

import (
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
	gwsAdapter "github.com/lxzan/uRouter/contrib/adapter/gws"
	httpAdapter "github.com/lxzan/uRouter/contrib/adapter/http"
	"github.com/lxzan/uRouter/contrib/codec/jsoniter"
	"github.com/lxzan/uRouter/contrib/log/zerolog"
	"net/http"
)

func init() {
	uRouter.SetLogger(zerolog.ZeroLogger)
	uRouter.SetJsonCodec(jsoniter.JsoniterCodec)
}

func main() {
	router := uRouter.New()
	router.Use(uRouter.AccessLog())

	upgrader := gws.NewUpgrader(func(c *gws.Upgrader) {
		c.EventHandler = &WebSocketHandler{adapter: gwsAdapter.NewAdapter(router)}
	})

	router.OnAction(http.MethodGet, "/connect", func(ctx *uRouter.Context) {
		socket, err := upgrader.Accept(ctx.Writer.Raw().(http.ResponseWriter), ctx.Request.Raw.(*http.Request))
		if err != nil {
			uRouter.Logger().Error(err.Error())
			return
		}
		go socket.Listen()
	})

	router.On("/greet", func(ctx *uRouter.Context) {
		_ = ctx.WriteString(int(gws.OpcodeText), "hello!")
	})

	router.Display()
	if err := http.ListenAndServe(":3000", httpAdapter.NewAdapter(router)); err != nil {
		uRouter.Logger().Panic(err.Error())
	}
}

type WebSocketHandler struct {
	gws.BuiltinEventEngine
	adapter *gwsAdapter.Adapter
}

func (c *WebSocketHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	if err := c.adapter.ServeWebSocket(socket, message); err != nil {
		uRouter.Logger().Error(err.Error())
	}
}

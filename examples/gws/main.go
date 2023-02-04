package main

import (
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
	ga "github.com/lxzan/uRouter/contrib/adapter/gws"
	"net/http"
)

func main() {
	var router = uRouter.New()

	router.On("/test", func(ctx *uRouter.Context) {
		ctx.WriteJSON(int(gws.OpcodeBinary), uRouter.A{
			"hello": "world!",
		})
	})

	var adapter = ga.NewAdapter(router, uRouter.BinaryHeader)
	var upgrader = gws.NewUpgrader(func(c *gws.Upgrader) {
		c.EventHandler = &WebSocket{adapter: adapter}
	})

	_ = http.ListenAndServe(":3000", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		socket, err := upgrader.Accept(writer, request)
		if err != nil {
			return
		}
		socket.Listen()
	}))
}

type WebSocket struct {
	adapter *ga.Adapter
}

func (c *WebSocket) OnOpen(socket *gws.Conn) {

}

func (c *WebSocket) OnError(socket *gws.Conn, err error) {
	println("onerror")
}

func (c *WebSocket) OnClose(socket *gws.Conn, code uint16, reason []byte) {
	println("onclose")
}

func (c *WebSocket) OnPing(socket *gws.Conn, payload []byte) {

}

func (c *WebSocket) OnPong(socket *gws.Conn, payload []byte) {

}

func (c *WebSocket) OnMessage(socket *gws.Conn, message *gws.Message) {
	if err := c.adapter.ServeWebSocket(socket, message.Bytes()); err != nil {
		return
	}
	message.Close()
}

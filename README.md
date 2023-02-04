# uRouter
universal router for http, websocket and other custom protocol

#### Feature
- use static router powered by map
- support middleware and router group 

#### Quick Start

- http
```go
package main

import (
	"github.com/lxzan/uRouter"
	httpAdapter "github.com/lxzan/uRouter/contrib/adapter/http"
	"net/http"
)

func main() {
	var router = uRouter.New()

	router.OnNoMatch = func(ctx *uRouter.Context) {
		ctx.WriteString(http.StatusNotFound, "not found")
	}

	group := router.Group("api/v1")

	group.On("greet", Greet)

	router.Display()

	http.ListenAndServe(":3000", httpAdapter.NewAdapter(router))
}

func Greet(ctx *uRouter.Context) {
	ctx.WriteJSON(http.StatusOK, uRouter.A{
		"hello": "world",
	})
}

```

- websocket (gws)
```go
// server
package main

import (
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
	gwsAdapter "github.com/lxzan/uRouter/contrib/adapter/gws"
	httpAdapter "github.com/lxzan/uRouter/contrib/adapter/http"
	"log"
	"net/http"
)

func main() {
	var router = uRouter.New()
	var upgrader = gws.NewUpgrader(func(c *gws.Upgrader) {
		c.EventHandler = &WebSocketHandler{
			adapter: gwsAdapter.NewAdapter(router, uRouter.TextHeader),
		}
	})

	router.On("/connect", func(ctx *uRouter.Context) {
		r, w := ctx.Raw()
		socket, err := upgrader.Accept(w.(http.ResponseWriter), r.(*http.Request))
		if err != nil {
			log.Println(err.Error())
			return
		}
		socket.Listen()
	})

	router.On("/greet", func(ctx *uRouter.Context) {
		var req = uRouter.A{}
		ctx.BindJSON(&req)
		log.Printf("%v", req)
	})

	if err := http.ListenAndServe(":3000", httpAdapter.NewAdapter(router)); err != nil {
		log.Fatalln(err.Error())
	}
}

type WebSocketHandler struct {
	gws.BuiltinEventEngine
	adapter *gwsAdapter.Adapter
}

func (c *WebSocketHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	if err := c.adapter.ServeWebSocket(socket, message.Bytes()); err != nil {
		log.Println(err.Error())
	}
}
```

```js
// client
let ws = new WebSocket('ws://127.0.0.1:3000/connect');
ws.send('0019{"X-Path":"/greet"}{"hello":"world"}');
```
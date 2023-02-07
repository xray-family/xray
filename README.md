# uRouter
universal router for http, websocket and custom protocol, one is all.

Hats off to express, koa, gin!

[![Build Status][1]][2] [![MIT licensed][3]][4] [![Go Version][5]][6] [![codecov][7]][8] [![Go Report Card][9]][10]

[1]: https://github.com/lxzan/uRouter/workflows/Go%20Test/badge.svg?branch=main

[2]: https://github.com/lxzan/uRouter/actions?query=branch%3Amain

[3]: https://img.shields.io/badge/license-MIT-blue.svg

[4]: LICENSE

[5]: https://img.shields.io/badge/go-%3E%3D1.16-30dff3?style=flat-square&logo=go

[6]: https://github.com/lxzan/uRouter

[7]: https://codecov.io/gh/lxzan/uRouter/branch/main/graph/badge.svg?token=0Tx9xH9Lvd

[8]: https://codecov.io/gh/lxzan/uRouter

[9]: https://goreportcard.com/badge/github.com/lxzan/uRouter

[10]: https://goreportcard.com/report/github.com/lxzan/uRouter

#### Feature
- no dependence
- dynamic separation matching request paths, powered by map and trie
- the onion model middleware, router group 
- adapt to `http`, `lxzan/gws`, `gorilla/websocket` ...

#### Index
- [uRouter](#urouter)
			- [Feature](#feature)
			- [Index](#index)
			- [Quick Start](#quick-start)
			- [Best Practice](#best-practice)
			- [Middleware](#middleware)


  
#### Quick Start

```go
package main

import (
	"github.com/lxzan/uRouter"
	http2 "github.com/lxzan/uRouter/contrib/adapter/http"
	"net/http"
)

func main() {
	var router = uRouter.New()
	router.Use(uRouter.Recovery(), uRouter.AccessLog())
	var group = router.Group("/api/v1")
	
	group.On("/user/:uid/article/:aid", func(ctx *uRouter.Context) {
		_ = ctx.WriteJSON(http.StatusOK, uRouter.Any{
			"uid": ctx.Param("uid"),
			"aid": ctx.Param("aid"),
		})
	})

	router.Display()
	_ = http.ListenAndServe(":3000", http2.NewAdapter(router))
}


```

#### Best Practice

- server
  
```go
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
	router.Use(uRouter.Recovery(), uRouter.AccessLog())

	var upgrader = gws.NewUpgrader(func(c *gws.Upgrader) {
		c.EventHandler = NewWebSocketHandler(router, gwsAdapter.NewAdapter(router))
	})

	router.On("/connect", func(ctx *uRouter.Context) {
		r := ctx.Request.Raw
		w := ctx.Writer.Raw()
		socket, err := upgrader.Accept(w.( http.ResponseWriter), r.(*http.Request))
		if err != nil {
			log.Println(err.Error())
			return
		}
		socket.Listen()
	})

	router.Display()
	if err := http.ListenAndServe(":3000", httpAdapter.NewAdapter(router)); err != nil {
		log.Fatalln(err.Error())
	}
}

func NewWebSocketHandler(router *uRouter.Router, adapter *gwsAdapter.Adapter) *WebSocketHandler {
	h := &WebSocketHandler{
		adapter: adapter,
		router:  router,
	}
	group := router.Group("", uRouter.Protocol(uRouter.ProtocolWebSocket))
	group.On("/greet", h.Greet)
	return h
}

type WebSocketHandler struct {
	gws.BuiltinEventEngine
	adapter *gwsAdapter.Adapter
	router  *uRouter.Router
}

func (c *WebSocketHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	if err := c.adapter.ServeWebSocket(socket, message); err != nil {
		log.Println(err.Error())
	}
}

func (c *WebSocketHandler) Greet(ctx *uRouter.Context) {
	var req = uRouter.Any{}
	if err := ctx.BindJSON(&req); err != nil {
		log.Println(err.Error())
		return
	}
	ctx.Writer.Header().Set(uRouter.ContentType, uRouter.MimeJson)
	_ = ctx.WriteJSON(int(gws.OpcodeText), req)
}
```

- client
  
```js
let ws = new WebSocket('ws://127.0.0.1:3000/connect');
ws.send('0019{"X-Path":"/greet"}{"hello":"world!"}');
```

#### Middleware
  
![onion model](https://upload-images.jianshu.io/upload_images/26203625-b80a51afcf265c9d.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/1078/format/webp)
  
```go
package main

import (
	"github.com/lxzan/uRouter"
	http2 "github.com/lxzan/uRouter/contrib/adapter/http"
	"net/http"
)

func main() {
	var router = uRouter.New()
	
	var list []int
	router.Use(func(ctx *uRouter.Context) {
		list = append(list, 1)
		ctx.Next()
		list = append(list, 2)
		fmt.Printf("%v\n", list)
	})
	
	var group = router.Group("/api/v1", func(ctx *uRouter.Context) {
		list = append(list, 3)
		ctx.Next()
		list = append(list, 4)
	})

	group.On("/greet", func(ctx *uRouter.Context) {
		list = append(list, 5)
	})

	router.Display()
	_ = http.ListenAndServe(":3000", http2.NewAdapter(router))
}
```

```
output: 1, 3, 5, 4, 2
```

# uRouter
universal router for http, websocket and other custom protocol

#### Feature
- use static router powered by map
- support middleware and router group 

#### Quick Start

##### Standard HTTP
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
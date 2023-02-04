package main

import (
	"github.com/lxzan/uRouter"
	ha "github.com/lxzan/uRouter/contrib/adapter/http"
	"net/http"
)

func main() {
	var router = uRouter.New()
	var ctl = new(Controller)
	ctl.Route(router)

	router.OnNoMatch = func(ctx *uRouter.Context) {
		ctx.Writer.Code(http.StatusNotFound)
		ctx.Writer.Write([]byte("not found"))
	}

	http.ListenAndServe(":3000", ha.NewAdapter(router))
}

type Controller struct {
}

func (c *Controller) Route(router *uRouter.Router) {
	router.On("/user/list", func(ctx *uRouter.Context) {
		ctx.WriteJSON(http.StatusNotFound, uRouter.A{
			"hello": "world",
		})
	})
}

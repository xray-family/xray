package main

import (
	"github.com/lxzan/uRouter"
	httpAdapter "github.com/lxzan/uRouter/contrib/adapter/http"
	"net/http"
)

type Controller struct {
	httpAdapter.RootController
}

func NewController() *Controller {
	return new(Controller)
}

func (c *Controller) Mapping(group *uRouter.Group) *Controller {
	group.OnEvent(http.MethodGet, "/greet/:name", c.Greet)
	group.OnPOST("/test", c.Test)
	return c
}

func (c *Controller) Greet(ctx *uRouter.Context) {
	_ = ctx.WriteJSON(http.StatusOK, uRouter.Any{
		"hello": ctx.Param("name"),
	})
}

func (c *Controller) Test(ctx *uRouter.Context) {
	_ = ctx.WriteString(http.StatusOK, "aha")
}

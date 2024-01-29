package main

import (
	"github.com/xray-family/xray"
	httpAdapter "github.com/xray-family/xray/contrib/adapter/http"
	"net/http"
)

type Controller struct {
	httpAdapter.RootController
}

func NewController() *Controller {
	return new(Controller)
}

func (c *Controller) Mapping(group *xray.Group) *Controller {
	group.OnEvent(http.MethodGet, "/greet/:name", c.Greet)
	group.POST("/test", c.Test)
	return c
}

func (c *Controller) Greet(ctx *xray.Context) {
	_ = ctx.WriteJSON(http.StatusOK, xray.Any{
		"hello": ctx.Param("name"),
	})
}

func (c *Controller) Test(ctx *xray.Context) {
	_ = ctx.WriteString(http.StatusOK, "aha")
}

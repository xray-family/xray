package uRouter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouteTree_Get(t *testing.T) {
	var as = assert.New(t)

	var tree = newRouteTree()
	tree.Set("/api/v1/user/:id", []HandlerFunc{AccessLog()})
	tree.Set("/api/v1/user/:id/profile", []HandlerFunc{AccessLog()})
	tree.Set("/api/v1/user/:id/article/:article_id", []HandlerFunc{AccessLog()})

	var list []string
	tree.Range(func(h *apiHandler) {
		list = append(list, h.VPath)
	})
	as.ElementsMatch(
		[]string{
			"/api/v1/user/:id",
			"/api/v1/user/:id/profile",
			"/api/v1/user/:id/article/:article_id",
		},
		list,
	)

	{
		handler := tree.Get("/user")
		as.Nil(handler)
	}
	{
		handler := tree.Get("")
		as.Nil(handler)
	}
	{
		handler := tree.Get("/api/v1/user/1")
		as.Equal(handler.VPath, "/api/v1/user/:id")
	}
	{
		handler := tree.Get("/api/v1/user/:id2")
		as.Equal(handler.VPath, "/api/v1/user/:id")
	}
	{
		handler := tree.Get("/api/v1/user/1/profile")
		as.Equal(handler.VPath, "/api/v1/user/:id/profile")
	}
	{
		handler := tree.Get("/api/v1/user/1/article/2")
		as.Equal(handler.VPath, "/api/v1/user/:id/article/:article_id")
	}
}

func TestRouteTree_Set(t *testing.T) {
	var tree = newRouteTree()
	tree.Set("", []HandlerFunc{AccessLog()})
	tree.Set("/api/", []HandlerFunc{AccessLog()})
}

func TestRouteTree_Range(t *testing.T) {
	var tree *routeTree
	tree.Range(func(h *apiHandler) {})
}

package uRouter

import (
	"github.com/lxzan/uRouter/internal"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRouteTree_Get(t *testing.T) {
	var as = assert.New(t)

	var tree = newRouteTree()
	tree.Set(&apiHandler{FullPath: "/api/v1/user/:id", Funcs: []HandlerFunc{AccessLog()}})
	tree.Set(&apiHandler{FullPath: "/api/v1/user/:id/profile", Funcs: []HandlerFunc{AccessLog()}})
	tree.Set(&apiHandler{FullPath: "/api/v1/user/:id/article/:article_id", Funcs: []HandlerFunc{AccessLog()}})

	var list []string
	tree.Range(func(h *apiHandler) {
		list = append(list, h.FullPath)
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
		handler, _ := tree.Get("/user")
		as.Nil(handler)
	}
	{
		handler, _ := tree.Get("")
		as.Nil(handler)
	}
	{
		handler, _ := tree.Get("/api/v1/user/1")
		as.Equal(handler.FullPath, "/api/v1/user/:id")
	}
	{
		handler, _ := tree.Get("/api/v1/user/:id2")
		as.Equal(handler.FullPath, "/api/v1/user/:id")
	}
	{
		handler, _ := tree.Get("/api/v1/user/1/profile")
		as.Equal(handler.FullPath, "/api/v1/user/:id/profile")
	}
	{
		handler, _ := tree.Get("/api/v1/user/1/article/2")
		as.Equal(handler.FullPath, "/api/v1/user/:id/article/:article_id")
	}
}

func TestRouteTree_Set(t *testing.T) {
	var tree = newRouteTree()
	tree.Set(&apiHandler{FullPath: "", Funcs: []HandlerFunc{AccessLog()}})
	tree.Set(&apiHandler{FullPath: "/api/", Funcs: []HandlerFunc{AccessLog()}})
}

func TestRouteTree_Range(t *testing.T) {
	var tree *routeTree
	tree.Range(func(h *apiHandler) {})
}

func BenchmarkRouteTree_Get(b *testing.B) {
	var count = 1024
	var segmentLen = 2
	var tree = newRouteTree()
	var r = internal.Numeric
	for i := 0; i < count; i++ {
		var idx = r.Intn(4)
		var list []string
		for j := 0; j < 4; j++ {
			var ele = string(r.Generate(segmentLen))
			if j == idx {
				ele = ":" + ele
			}
			list = append(list, ele)
		}

		api := &apiHandler{FullPath: strings.Join(list, SEP)}
		tree.Set(api)
	}

	var paths []string
	for i := 0; i < count; i++ {
		var path = r.Generate(12)
		path[0], path[3], path[6], path[9] = SEP[0], SEP[0], SEP[0], SEP[0]
		paths = append(paths, string(path))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var path = paths[i&(count-1)]
		tree.Get(path)
	}
}

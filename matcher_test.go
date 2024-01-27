package xray

import (
	"github.com/lxzan/xray/internal"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRouteTree_Get(t *testing.T) {
	var as = assert.New(t)

	var tree = newRouteTree()
	tree.Set(&apiHandler{Path: "/api/v1/user/:id", Funcs: []HandlerFunc{AccessLog()}})
	tree.Set(&apiHandler{Path: "/api/v1/user/:id/profile", Funcs: []HandlerFunc{AccessLog()}})
	tree.Set(&apiHandler{Path: "/api/v1/user/:id/article/:article_id", Funcs: []HandlerFunc{AccessLog()}})

	var list []string
	tree.Range(func(h *apiHandler) {
		list = append(list, h.Path)
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
		handler, _ := tree.Get("/api/v1/user/1")
		as.Equal(handler.Path, "/api/v1/user/:id")
	}
	{
		handler, _ := tree.Get("/user")
		as.Nil(handler)
	}
	{
		handler, _ := tree.Get("")
		as.Nil(handler)
	}
	{
		handler, _ := tree.Get("/api/v1/user/:id2")
		as.Equal(handler.Path, "/api/v1/user/:id")
	}
	{
		handler, _ := tree.Get("/api/v1/user/1/profile")
		as.Equal(handler.Path, "/api/v1/user/:id/profile")
	}
	{
		handler, _ := tree.Get("/api/v1/user/1/article/2")
		as.Equal(handler.Path, "/api/v1/user/:id/article/:article_id")
	}
}

func TestRouteTree_Set(t *testing.T) {
	var tree = newRouteTree()
	tree.Set(&apiHandler{Path: "", Funcs: []HandlerFunc{AccessLog()}})
	tree.Set(&apiHandler{Path: "/api/", Funcs: []HandlerFunc{AccessLog()}})
}

func TestRouteTree_Range(t *testing.T) {
	var tree *routeTree
	tree.Range(func(h *apiHandler) {})
}

func BenchmarkRouteTree_Get(b *testing.B) {
	var count = 1024
	var segmentLen = 4
	var tree = newRouteTree()
	var r = internal.AlphabetNumeric
	var list []string
	for i := 0; i < count; i++ {
		var idx = r.Intn(4)
		var arr []string
		for j := 0; j < 4; j++ {
			var ele = string(r.Generate(segmentLen))
			if j == idx {
				ele = ":" + ele
			}
			arr = append(arr, ele)
		}

		path := SEP + strings.Join(arr, SEP)
		list = append(list, path)
		api := &apiHandler{Path: path}
		tree.Set(api)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var path = list[i&(count-1)]
		tree.Get(path)
	}
}

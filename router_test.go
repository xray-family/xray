package uRouter

import (
	"github.com/lxzan/uRouter/constant"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	var as = assert.New(t)

	t.Run("static router", func(t *testing.T) {
		var list []int
		var r = New()
		r.Use(func(ctx *Context) {
			list = append(list, 1)
			ctx.Next()
			list = append(list, 2)
		})

		var g0 = r.Group("", func(ctx *Context) {
			list = append(list, 3)
			ctx.Next()
			list = append(list, 4)
		})

		var g1 = g0.Group("/api/v1", func(ctx *Context) {
			list = append(list, 5)
			ctx.Next()
			list = append(list, 6)
		})

		g1.OnEvent(http.MethodGet, "greet", func(ctx *Context) {
			list = append(list, 9)
		}, func(ctx *Context) {
			list = append(list, 7)
			ctx.Next()
			list = append(list, 8)
		})

		r.doStart()

		path := "/api/v1/greet"
		ctx := NewContext(&Request{}, newResponseWriterMocker())
		r.EmitEvent(http.MethodGet, path, ctx)

		as.Equal(9, len(list))
		as.Equal(1, list[0])
		as.Equal(3, list[1])
		as.Equal(5, list[2])
		as.Equal(7, list[3])
		as.Equal(9, list[4])
		as.Equal(8, list[5])
		as.Equal(6, list[6])
		as.Equal(4, list[7])
		as.Equal(2, list[8])
	})

	t.Run("dynamic router", func(t *testing.T) {
		var list []int
		var r = New()
		r.Use(func(ctx *Context) {
			list = append(list, 1)
			ctx.Next()
			list = append(list, 2)
		})

		var g0 = r.Group("", func(ctx *Context) {
			list = append(list, 3)
			ctx.Next()
			list = append(list, 4)
		})

		var g1 = g0.Group("/api/v1", func(ctx *Context) {
			list = append(list, 5)
			ctx.Next()
			list = append(list, 6)
		})

		g1.On("greet/:id", func(ctx *Context) {
			list = append(list, 9)
			as.Equal("1", ctx.Param("id"))
		}, func(ctx *Context) {
			list = append(list, 7)
			ctx.Next()
			list = append(list, 8)
		})

		path := "/api/v1/greet/1"
		ctx := NewContext(
			&Request{Header: NewHttpHeader(http.Header{"X-Path": []string{path}})},
			newResponseWriterMocker(),
		)
		r.doStart()

		r.Emit(path, ctx)
		as.Equal(9, len(list))
		as.Equal(1, list[0])
		as.Equal(3, list[1])
		as.Equal(5, list[2])
		as.Equal(7, list[3])
		as.Equal(9, list[4])
		as.Equal(8, list[5])
		as.Equal(6, list[6])
		as.Equal(4, list[7])
		as.Equal(2, list[8])
	})

	t.Run("mix router", func(t *testing.T) {
		var list []int
		var r = New()
		r.Use(func(ctx *Context) {
			list = append(list, 1)
			ctx.Next()
			list = append(list, 2)
		})

		var g0 = r.Group("", func(ctx *Context) {
			list = append(list, 3)
			ctx.Next()
			list = append(list, 4)
		})

		var g1 = g0.Group("/api/:version", func(ctx *Context) {
			list = append(list, 5)
			ctx.Next()
			list = append(list, 6)
		})

		g1.On("greet", func(ctx *Context) {
			list = append(list, 9)
			as.Equal("v1", ctx.Param("version"))
		}, func(ctx *Context) {
			list = append(list, 7)
			ctx.Next()
			list = append(list, 8)
		})

		r.doStart()

		path := "/api/v1/greet"
		ctx := NewContext(
			&Request{Header: NewHttpHeader(http.Header{"X-Path": []string{path}})},
			newResponseWriterMocker(),
		)
		r.Emit(path, ctx)

		as.Equal(9, len(list))
		as.Equal(1, list[0])
		as.Equal(3, list[1])
		as.Equal(5, list[2])
		as.Equal(7, list[3])
		as.Equal(9, list[4])
		as.Equal(8, list[5])
		as.Equal(6, list[6])
		as.Equal(4, list[7])
		as.Equal(2, list[8])
	})

	t.Run("match", func(t *testing.T) {
		var r = New()
		var list []int

		r.On("test", func(ctx *Context) {
			list = append(list, 3)
		}, func(ctx *Context) {
			list = append(list, 1)
			ctx.Next()
			list = append(list, 2)
		})
		r.doStart()

		path := "/test"
		ctx := NewContext(
			&Request{Header: NewHttpHeader(http.Header{"X-Path": []string{path}}), Body: nil},
			newResponseWriterMocker(),
		)
		r.Emit(path, ctx)

		as.Equal(len(list), 3)
		as.Equal(1, list[0])
		as.Equal(3, list[1])
		as.Equal(2, list[2])
	})

	t.Run("no match", func(t *testing.T) {
		var r = New()
		var list []int

		r.OnNotFound = func(ctx *Context) { list = append(list, 1) }
		r.doStart()

		path := "/test"
		ctx := NewContext(
			&Request{Header: NewHttpHeader(http.Header{"X-Path": []string{}}), Body: nil},
			newResponseWriterMocker(),
		)
		r.Emit(path, ctx)

		as.Equal(len(list), 1)
		as.Equal(1, list[0])
	})

	t.Run("no handler", func(t *testing.T) {
		var r = New()
		var list []int

		r.Group("test")

		path := "/test"
		ctx := NewContext(
			&Request{
				Header: NewHttpHeader(http.Header{constant.XPath: []string{path}}), Body: nil,
			},
			newResponseWriterMocker(),
		)
		r.Emit(path, ctx)

		r.staticRoutes["404"] = nil

		as.Equal(len(list), 0)
	})
}

func TestRouter_OnNoMatch(t *testing.T) {
	var as = assert.New(t)

	t.Run("", func(t *testing.T) {
		var r = New()
		r.Use(func(ctx *Context) {
			ctx.Set("sum", 1)
			ctx.Next()
		}, func(ctx *Context) {
			val, _ := ctx.Get("sum")
			ctx.Set("sum", val.(int)+2)
			ctx.Next()
		})
		r.doStart()

		const count = 10
		var wg = &sync.WaitGroup{}
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func() {
				var path = "test"
				var ctx = NewContext(
					&Request{Header: NewHttpHeader(http.Header{constant.XPath: []string{path}})},
					newResponseWriterMocker(),
				)
				r.Emit(path, ctx)
				sum, _ := ctx.Get("sum")
				as.Equal(3, sum.(int))
				wg.Done()
			}()
		}
		wg.Wait()
	})

	t.Run("static router", func(t *testing.T) {
		var list []int
		var r = New()
		r.Use(func(ctx *Context) {
			list = append(list, 1)
			ctx.Next()
			list = append(list, 2)
		})

		var g0 = r.Group("", func(ctx *Context) {
			list = append(list, 3)
			ctx.Next()
			list = append(list, 4)
		})

		var g1 = g0.Group("/api/v1", func(ctx *Context) {
			list = append(list, 5)
			ctx.Next()
			list = append(list, 6)
		})

		g1.On("greet", func(ctx *Context) {
			list = append(list, 9)
		}, func(ctx *Context) {
			list = append(list, 7)
			ctx.Next()
			list = append(list, 8)
		})

		r.OnNotFound = func(ctx *Context) { list = append(list, 10) }
		r.doStart()

		path := "/api/v1/xxx"
		ctx := NewContext(
			&Request{Header: NewHttpHeader(http.Header{"X-Path": []string{path}})},
			newResponseWriterMocker(),
		)
		r.Emit(path, ctx)

		as.Equal(3, len(list))
		as.Equal(1, list[0])
		as.Equal(10, list[1])
		as.Equal(2, list[2])
	})
}

func TestRouter_Conflict(t *testing.T) {
	var as = assert.New(t)

	t.Run("route conflict 1", func(t *testing.T) {
		defer func() {
			e := recover()
			as.NotNil(e)
		}()

		var r = New()
		var g = r.Group("user")
		g.On("1", AccessLog())
		r.On("user/1", AccessLog())
		r.doStart()
	})

	t.Run("route conflict 2", func(t *testing.T) {
		defer func() {
			e := recover()
			as.NotNil(e)
		}()

		var r = New()
		r.On("user/:id", AccessLog())
		r.On("user/1", AccessLog())
		r.doStart()
	})

	t.Run("route conflict 3", func(t *testing.T) {
		defer func() {
			e := recover()
			as.NotNil(e)
		}()

		var r = New()
		r.On("user/:id", AccessLog())
		var g = r.Group("user")
		g.On("1", AccessLog())
		r.doStart()
	})

	t.Run("route conflict 4", func(t *testing.T) {
		defer func() {
			e := recover()
			as.NotNil(e)
		}()

		var r = New()
		r.On("user/1", AccessLog())
		var g = r.Group("user")
		g.On(":id", AccessLog())
		r.doStart()
	})

	t.Run("route conflict 5", func(t *testing.T) {
		defer func() {
			e := recover()
			as.NotNil(e)
		}()

		var r = New()
		r.On("user/:id", AccessLog())
		r.On("user/:name", AccessLog())
		r.doStart()
	})
}

func TestRouter_Display(t *testing.T) {
	r := New()
	r.OnEvent(http.MethodGet, "/user/list", func(ctx *Context) {})
	r.OnEvent(http.MethodPost, "/user/:id", func(ctx *Context) {})
	r.doStart()
}

func TestRouter_Dynamic(t *testing.T) {
	var as = assert.New(t)

	t.Run("", func(t *testing.T) {
		defer func() {
			e := recover()
			as.NotNil(e)
		}()
		r := New()
		r.OnEvent(http.MethodGet, "/user/list", func(ctx *Context) {})
		r.OnEvent(http.MethodGet, "/user/:id", func(ctx *Context) {})
		r.doStart()
	})

	t.Run("", func(t *testing.T) {
		defer func() {
			e := recover()
			as.Nil(e)
		}()
		r := New()
		r.OnEvent(http.MethodGet, "/user/list", func(ctx *Context) {})
		r.OnEvent(http.MethodDelete, "/user/:id", func(ctx *Context) {})
		r.OnEvent(http.MethodPost, "/user/:id", func(ctx *Context) {})
	})

	t.Run("", func(t *testing.T) {
		defer func() {
			e := recover()
			as.Nil(e)
		}()
		r := New()

		sum := 0
		r.OnEvent(http.MethodGet, "/user/:id/profile", func(ctx *Context) {
			sum++
		})
		r.OnEvent(http.MethodGet, "/user/:id", func(ctx *Context) {

		})
		r.doStart()

		ctx := NewContext(&Request{}, newResponseWriterMocker())
		r.EmitEvent(http.MethodGet, "/user/1/profile", ctx)
		as.Equal(1, sum)
	})

	t.Run("", func(t *testing.T) {
		defer func() {
			e := recover()
			as.Nil(e)
		}()
		r := New()

		sum := 0
		r.OnEvent(http.MethodGet, "/user/:id/profile", func(ctx *Context) {
			sum++
		})
		r.OnEvent(http.MethodGet, "/user/:id", func(ctx *Context) {

		})
		ctx := NewContext(&Request{}, newResponseWriterMocker())
		r.EmitEvent(http.MethodPost, "/user/1/profile", ctx)
		as.Equal(0, sum)
	})
}

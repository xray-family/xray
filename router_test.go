package uRouter

import (
	"github.com/lxzan/uRouter/internal"
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

		r.StartSilently()

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
		r.StartSilently()

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

		r.StartSilently()

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
		r.StartSilently()

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
		r.StartSilently()

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
				Header: NewHttpHeader(http.Header{UPath: []string{path}}), Body: nil,
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
		r.StartSilently()

		const count = 10
		var wg = &sync.WaitGroup{}
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func() {
				var path = "test"
				var ctx = NewContext(
					&Request{Header: NewHttpHeader(http.Header{UPath: []string{path}})},
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
		r.StartSilently()

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
		r.StartSilently()
	})

	t.Run("route conflict 2", func(t *testing.T) {
		defer func() {
			e := recover()
			as.NotNil(e)
		}()

		var r = New()
		r.On("user/:id", AccessLog())
		r.On("user/1", AccessLog())
		r.StartSilently()
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
		r.StartSilently()
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
		r.StartSilently()
	})

	t.Run("route conflict 5", func(t *testing.T) {
		defer func() {
			e := recover()
			as.NotNil(e)
		}()

		var r = New()
		r.On("user/:id", AccessLog())
		r.On("user/:name", AccessLog())
		r.StartSilently()
	})
}

func TestRouter_Display(t *testing.T) {
	r := New()
	r.OnEvent(http.MethodGet, "/user/list", func(ctx *Context) {})
	r.OnEvent(http.MethodPost, "/user/:id", func(ctx *Context) {})
	r.StartSilently()
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
		r.StartSilently()
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
		r.StartSilently()

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

func TestRouter_EmitRandom(t *testing.T) {
	var as = assert.New(t)

	var count = 1024
	var paths []string
	for i := 0; i < count; i++ {
		s0 := string(internal.AlphabetNumeric.Generate(4))
		s1 := string(internal.AlphabetNumeric.Generate(4))
		s2 := string(internal.AlphabetNumeric.Generate(4))
		s3 := string(internal.AlphabetNumeric.Generate(4))

		flag := internal.AlphabetNumeric.Intn(32)
		switch flag {
		case 1:
			s0 = ":" + s0
		case 3:
			s1 = ":" + s1
		case 5:
			s2 = ":" + s2
		case 7:
			s3 = ":" + s3
		}

		paths = append(paths, internal.JoinPath(SEP, s0, s1, s2, s3))
	}

	var mapping = map[string]uint8{}
	var r = New()
	var g0 = r.Group("api")
	var g1 = g0.Group("v1")
	for _, v := range paths {
		g1.OnGET(v, func(ctx *Context) {
			mapping[ctx.Request.VPath] = 1
		})
	}
	r.StartSilently()

	var prefix = "/api/v1"
	for i := 0; i < count; i++ {
		ctx := newContextMocker()
		r.EmitEvent(http.MethodGet, prefix+paths[i], ctx)
	}
	as.Equal(count, len(mapping))

	var exists = func(p string) bool {
		for _, v := range r.apis {
			if v.Path == p {
				return true
			}

			list1 := internal.Split(v.Path, SEP)
			list2 := internal.Split(p, SEP)
			if len(list1) != len(list2) {
				continue
			}

			var counter = 0
			for i, seg := range list1 {
				if seg == list2[i] || isVar(seg) {
					counter++
				}
			}
			if counter == len(list1) {
				return true
			}
		}
		return false
	}

	//遍历去验证
	var expected = 0
	for k, _ := range mapping {
		delete(mapping, k)
	}
	for i := 0; i < count; i++ {
		ctx := newContextMocker()
		segments := internal.Split(paths[i], SEP)
		flag := internal.AlphabetNumeric.Intn(32)
		switch flag {
		case 2:
			segments[0] = string(internal.Numeric.Generate(4))
		case 4:
			segments[1] = string(internal.Numeric.Generate(4))
		case 6:
			segments[2] = string(internal.Numeric.Generate(4))
		case 8:
			segments[3] = string(internal.Numeric.Generate(4))
		}
		var arr = []string{prefix}
		arr = append(arr, segments...)
		path := internal.JoinPath(SEP, arr...)
		if exists(path) {
			expected++
		}
		r.EmitEvent(http.MethodGet, path, ctx)
	}
	as.Equal(expected, len(mapping))
}

func TestRouter_Actions(t *testing.T) {
	var as = assert.New(t)

	t.Run("get", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.OnGET("/test", func(ctx *Context) {
			sum++
		})
		r.StartSilently()
		r.EmitEvent(http.MethodGet, "/test", newContextMocker())
		as.Equal(1, sum)
	})

	t.Run("post", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.OnPOST("/test", func(ctx *Context) {
			sum++
		})
		r.StartSilently()
		r.EmitEvent(http.MethodPost, "/test", newContextMocker())
		as.Equal(1, sum)
	})

	t.Run("put", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.OnPUT("/test", func(ctx *Context) {
			sum++
		})
		r.StartSilently()
		r.EmitEvent(http.MethodPut, "/test", newContextMocker())
		as.Equal(1, sum)
	})

	t.Run("delete", func(t *testing.T) {
		var r = New()
		var sum = 0
		r.OnDELETE("/test", func(ctx *Context) {
			sum++
		})
		r.StartSilently()
		r.EmitEvent(http.MethodDelete, "/test", newContextMocker())
		as.Equal(1, sum)
	})
}

func TestGroup_Actions(t *testing.T) {
	var as = assert.New(t)

	t.Run("get", func(t *testing.T) {
		var r = New()
		var g = r.Group("api/v1")
		var sum = 0
		g.OnGET("/test", func(ctx *Context) {
			sum++
		})
		r.StartSilently()
		r.EmitEvent(http.MethodGet, "/api/v1/test", newContextMocker())
		as.Equal(1, sum)
	})

	t.Run("post", func(t *testing.T) {
		var r = New()
		var g = r.Group("api/v1")
		var sum = 0
		g.OnPOST("/test", func(ctx *Context) {
			sum++
		})
		r.StartSilently()
		r.EmitEvent(http.MethodPost, "/api/v1/test", newContextMocker())
		as.Equal(1, sum)
	})

	t.Run("put", func(t *testing.T) {
		var r = New()
		var g = r.Group("api/v1")
		var sum = 0
		g.OnPUT("/test", func(ctx *Context) {
			sum++
		})
		r.StartSilently()
		r.EmitEvent(http.MethodPut, "/api/v1/test", newContextMocker())
		as.Equal(1, sum)
	})

	t.Run("delete", func(t *testing.T) {
		var r = New()
		var g = r.Group("api/v1")
		var sum = 0
		g.OnDELETE("/test", func(ctx *Context) {
			sum++
		})
		r.StartSilently()
		r.EmitEvent(http.MethodDelete, "/api/v1/test", newContextMocker())
		as.Equal(1, sum)
	})
}

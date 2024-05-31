package treemap

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestTreeMap(t *testing.T) {
	var tm = New[uint8]()
	tm.Set(http.MethodGet, "/api/v1/test", 1)
	tm.Set(http.MethodGet, "/api/greet/:name/id", 3)
	tm.Set(http.MethodGet, "/api/greet/:name", 2)

	v, _ := tm.Get(http.MethodGet, "/api/v1/test")
	assert.Equal(t, v, uint8(1))

	v, _ = tm.Get(http.MethodGet, "/api/greet/lee")
	assert.Equal(t, v, uint8(2))

	v, _ = tm.Get(http.MethodGet, "/api/greet/wang")
	assert.Equal(t, v, uint8(2))

	v, _ = tm.Get(http.MethodGet, "/api/greet/wang/id")
	assert.Equal(t, v, uint8(3))
}

func TestTreeMap_Exists(t *testing.T) {
	var tm = New[uint8]()
	tm.Set(http.MethodGet, "/api/v1/test", 1)
	tm.Set(http.MethodGet, "/api/greet/:name", 2)
	tm.Set(http.MethodGet, "/api/greet/:name/id", 3)
	v, _ := tm.Exists(http.MethodGet, "/api/v1/:any")
	assert.Equal(t, v, uint8(1))

	v, _ = tm.Exists(http.MethodGet, "/api/v1/test")
	assert.Equal(t, v, uint8(1))

	v, _ = tm.Exists(http.MethodGet, "/api/greet/lee")
	assert.Equal(t, v, uint8(2))

	v, _ = tm.Exists(http.MethodGet, "/api/greet/lee/id")
	assert.Equal(t, v, uint8(3))

	v, _ = tm.Exists(http.MethodGet, "/api/greet/:name/:id")
	assert.Equal(t, v, uint8(3))
}

func TestTreeMap_Range(t *testing.T) {
	var tm = New[uint8]()
	tm.Set(http.MethodGet, "/api/v1/test", 1)
	tm.Set(http.MethodGet, "/api/greet/:name", 2)
	tm.Set(http.MethodGet, "/api/greet/:name/id", 3)
	var list []uint8
	tm.Range(func(h uint8) { list = append(list, h) })
	assert.ElementsMatch(t, list, []uint8{1, 2, 3})
}

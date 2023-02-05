package uRouter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAny(t *testing.T) {
	var as = assert.New(t)
	var m = Any{
		"num1": 1,
		"num2": int64(2),
		"str":  "3",
	}
	as.Equal(1, m.ToInt("num1"))
	as.Equal(int64(2), m.ToInt64("num2"))
	as.Equal("3", m.ToString("str"))

	as.Equal(false, m.Exists("xxx"))
	as.Equal(0, m.ToInt("xxx"))
	as.Equal(int64(0), m.ToInt64("xxx"))
	as.Equal("", m.ToString("xxx"))
}

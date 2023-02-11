package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoinPath(t *testing.T) {
	var as = assert.New(t)
	var sep = "/"
	as.Equal("/", JoinPath(sep))
	as.Equal("/", JoinPath(sep, ""))
	as.Equal("/api", JoinPath(sep, "api"))
	as.Equal("/api/v1", JoinPath(sep, "api", "v1"))
	as.Equal("/api/v1", JoinPath(sep, "api", "/v1"))
	as.Equal("/api/v1", JoinPath(sep, "api/", "v1"))
	as.Equal("/api/v1", JoinPath(sep, "//api", "//v1"))
	as.Equal("/api/v1", JoinPath(sep, "//api//", "//v1"))
	as.Equal("/api/v1", JoinPath(sep, "/api/", "/v1/"))
	as.Equal("/api/v1", JoinPath(sep, "//api//", "//v1//"))
	as.Equal("/api/v1", JoinPath(sep, "//api////v1//"))
	as.Equal("/api/v1", JoinPath(sep, "/api/v1/"))
	as.Equal("/api/v1", JoinPath(sep, "/api/v1"))
	as.Equal("/", JoinPath(sep, "/"))
	as.Equal("/api/v1/user/list", JoinPath(sep, "/api/v1", "user", "list"))
	as.Equal("/api/v1/user/list", JoinPath(sep, "/api/v1", "user", "list/"))
	as.Equal("/api/v1/user/list", JoinPath(sep, "/api/v1", "/user", "/list"))
}

func TestRandomString(t *testing.T) {
	var s = AlphabetNumeric.Generate(16)
	assert.Equal(t, 16, len(s))
	assert.Less(t, AlphabetNumeric.Intn(16), 16)
}

func BenchmarkJoinPath1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JoinPath("/", "/api/v1/user/list")
	}
}

func BenchmarkJoinPath2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JoinPath("/", "/api/v1", "user", "list")
	}
}

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
	as.Equal("/api/v1/user", JoinPath(sep, "/api/v1", "", "/user"))
	as.Equal("/api/v1/user", JoinPath(sep, "/api/v1//user"))
}

func TestTrimPath(t *testing.T) {
	var as = assert.New(t)
	as.Equal("/api/v1", TrimPath("api/v1//"))
	as.Equal("/", TrimPath(""))
	as.Equal("/", TrimPath("/"))
	as.Equal("/", TrimPath("//"))
	as.Equal("/api", TrimPath("//api/"))
	as.Equal("/api", TrimPath("//api//"))
	as.Equal("/api/v1", TrimPath("//api/v1//"))
	as.Equal("/api/v1", TrimPath("///api/v1//"))
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

func TestSelectString(t *testing.T) {
	assert.Equal(t, "1", SelectString(true, "1", "0"))
	assert.Equal(t, "0", SelectString(false, "1", "0"))
}

func TestSplit(t *testing.T) {
	assert.ElementsMatch(t, []string{"api", "v1"}, Split("/api/v1"))
	assert.ElementsMatch(t, []string{"api", "v1"}, Split("/api/v1/"))
	assert.ElementsMatch(t, []string{"ming", "hong", "hu"}, Split("ming/ hong/ hu"))
	assert.ElementsMatch(t, []string{"ming", "hong", "hu"}, Split("/ming/ hong/ hu/ "))
	assert.ElementsMatch(t, []string{"ming", "hong", "hu"}, Split("\nming/ hong/ hu\n"))
}

func TestGetMaxLength(t *testing.T) {
	var list = []string{"1", "12", "123"}
	assert.Equal(t, 3, GetMaxLength(list...))
}

func TestPadding(t *testing.T) {
	assert.Equal(t, 16, len(Padding("123", 16)))
}

func TestFastSplit(t *testing.T) {
	var as = assert.New(t)

	{
		var s = "/api/v1/list"
		var list []string
		FastSplit(s, func(segment string) bool {
			list = append(list, segment)
			return true
		})
		as.ElementsMatch([]string{"api", "v1", "list"}, list)
	}

	{
		var s = "/api//v1/list"
		var list []string
		FastSplit(s, func(segment string) bool {
			list = append(list, segment)
			return true
		})
		as.ElementsMatch([]string{"api", "", "v1", "list"}, list)
	}

	{
		var s = "/api/v1/user/list"
		var list []string
		FastSplit(s, func(segment string) bool {
			list = append(list, segment)
			return len(list) <= 2
		})
		as.ElementsMatch([]string{"api", "v1", "user"}, list)
	}

	{
		var s = "/api"
		var list []string
		FastSplit(s, func(segment string) bool {
			list = append(list, segment)
			return true
		})
		as.ElementsMatch([]string{"api"}, list)
	}

	{
		var s = "/"
		var list []string
		FastSplit(s, func(segment string) bool {
			list = append(list, segment)
			return true
		})
		as.ElementsMatch([]string{}, list)
	}

	{
		var s = ""
		var list []string
		FastSplit(s, func(segment string) bool {
			list = append(list, segment)
			return true
		})
		as.ElementsMatch([]string{}, list)
	}
}

func TestClone(t *testing.T) {
	var a = []int{1, 3, 5}
	var b = Clone(a)
	assert.ElementsMatch(t, a, b)
}

func TestCatch(t *testing.T) {
	var err = Catch(func() {
		panic(1)
	})
	assert.Error(t, err)
}

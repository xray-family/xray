package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoin1(t *testing.T) {
	var as = assert.New(t)
	var sep = "/"
	as.Equal("/", Join1("", sep))
	as.Equal("/api", Join1("api", sep))
	as.Equal("/api", Join1("/api", sep))
	as.Equal("/api", Join1("/api/", sep))
}

func TestJoin2(t *testing.T) {
	var as = assert.New(t)
	var sep = "/"
	as.Equal("/", Join2("", "", sep))
	as.Equal("/api", Join2("", "/api", sep))
	as.Equal("/api", Join2("/api/", "", sep))
	as.Equal("/api/v1", Join2("/api/", "v1", sep))
	as.Equal("/api/v1", Join2("/api", "v1", sep))
	as.Equal("/api/v1", Join2("/api/", "/v1", sep))
	as.Equal("/api/v1", Join2("/api/", "/v1/", sep))
}

func TestRandomString(t *testing.T) {
	var s = AlphabetNumeric.Generate(16)
	assert.Equal(t, 16, len(s))
	assert.Less(t, AlphabetNumeric.Intn(16), 16)
}

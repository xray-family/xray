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
}

func TestSelectString(t *testing.T) {
	assert.Equal(t, "1", SelectString(true, "1", "0"))
	assert.Equal(t, "0", SelectString(false, "1", "0"))
}

package codec

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestStdJsonCodec_Encode(t *testing.T) {
	var v = http.Header{}
	v.Set("content-type", "application/json")
	_, err := StdJsonCodec.Encode(v)
	assert.NoError(t, err)
}

func TestStdJsonCodec_Decode(t *testing.T) {
	var d = `{"XPath":["/api"]}`
	var v = http.Header{}
	err := StdJsonCodec.Decode([]byte(d), &v)
	assert.NoError(t, err)
}

func TestStdJsonCodec(t *testing.T) {
	StdJsonCodec.NewEncoder(bytes.NewBufferString(""))
	StdJsonCodec.NewDecoder(bytes.NewBufferString(""))
}

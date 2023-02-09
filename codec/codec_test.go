package codec

import (
	"github.com/lxzan/uRouter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStdJsonCodec_Encode(t *testing.T) {
	var v = uRouter.MapHeader{}
	v.Set(uRouter.ContentType, uRouter.MimeJson)
	_, err := uRouter.StdJsonCodec.Encode(v)
	assert.NoError(t, err)
}

func TestStdJsonCodec_Decode(t *testing.T) {
	var d = `{"xpath":"/api"}`
	var v = uRouter.MapHeader{}
	err := uRouter.StdJsonCodec.Decode([]byte(d), &v)
	assert.NoError(t, err)
}

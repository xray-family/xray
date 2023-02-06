package jsoniter

import (
	json "github.com/json-iterator/go"
	"github.com/lxzan/uRouter"
	"io"
)

var JsoniterCodec = &codec{caller: json.ConfigFastest}

type codec struct {
	caller json.API
}

func (c codec) NewEncoder(w io.Writer) uRouter.Encoder {
	return c.caller.NewEncoder(w)
}

func (c codec) NewDecoder(r io.Reader) uRouter.Decoder {
	return c.caller.NewDecoder(r)
}

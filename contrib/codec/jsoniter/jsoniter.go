package jsoniter

import (
	json "github.com/json-iterator/go"
	"github.com/lxzan/uRouter"
	"io"
)

var JsoniterCodec = &codec{caller: json.ConfigFastest}

func (c codec) Encode(v interface{}) ([]byte, error) {
	return c.caller.Marshal(v)
}

func (c codec) Decode(data []byte, v interface{}) error {
	return c.caller.Unmarshal(data, v)
}

type codec struct {
	caller json.API
}

func (c codec) NewEncoder(w io.Writer) uRouter.Encoder {
	return c.caller.NewEncoder(w)
}

func (c codec) NewDecoder(r io.Reader) uRouter.Decoder {
	return c.caller.NewDecoder(r)
}

package jsoniter

import (
	json "github.com/json-iterator/go"
	"github.com/lxzan/uRouter/codec"
	"io"
)

var JsoniterCodec = &jsonCodec{caller: json.ConfigFastest}

type jsonCodec struct {
	caller json.API
}

func (c jsonCodec) Encode(v interface{}) ([]byte, error) {
	return c.caller.Marshal(v)
}

func (c jsonCodec) Decode(data []byte, v interface{}) error {
	return c.caller.Unmarshal(data, v)
}

func (c jsonCodec) NewEncoder(w io.Writer) codec.Encoder {
	return c.caller.NewEncoder(w)
}

func (c jsonCodec) NewDecoder(r io.Reader) codec.Decoder {
	return c.caller.NewDecoder(r)
}

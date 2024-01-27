package jsoniter

import (
	json "github.com/json-iterator/go"
	"github.com/lxzan/xray/codec"
	"io"
)

var JsoniterCodec = &Codec{API: json.ConfigFastest}

type Codec struct {
	API json.API
}

func (c Codec) Encode(v interface{}) ([]byte, error) {
	return c.API.Marshal(v)
}

func (c Codec) Decode(data []byte, v interface{}) error {
	return c.API.Unmarshal(data, v)
}

func (c Codec) NewEncoder(w io.Writer) codec.Encoder {
	return c.API.NewEncoder(w)
}

func (c Codec) NewDecoder(r io.Reader) codec.Decoder {
	return c.API.NewDecoder(r)
}

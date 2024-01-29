package jsoniter

import (
	json "github.com/json-iterator/go"
	"github.com/lxzan/xray/codec"
	"io"
)

var Codec codec.Codec = &JsonCodec{API: json.ConfigFastest}

type JsonCodec struct{ API json.API }

func (c *JsonCodec) Encode(v any) ([]byte, error) { return c.API.Marshal(v) }

func (c *JsonCodec) Decode(data []byte, v any) error { return c.API.Unmarshal(data, v) }

func (c *JsonCodec) NewEncoder(w io.Writer) codec.Encoder { return c.API.NewEncoder(w) }

func (c *JsonCodec) NewDecoder(r io.Reader) codec.Decoder { return c.API.NewDecoder(r) }

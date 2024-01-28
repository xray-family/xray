package codec

import (
	"encoding/json"
	"io"
)

var StdJsonCodec = new(stdJsonCodec)

type (
	Encoder interface {
		Encode(v any) error
	}

	Decoder interface {
		Decode(v any) error
	}

	Codec interface {
		NewEncoder(w io.Writer) Encoder
		NewDecoder(r io.Reader) Decoder
		Encode(v any) ([]byte, error)
		Decode(data []byte, v any) error
	}
)

type stdJsonCodec struct{}

func (c *stdJsonCodec) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (c *stdJsonCodec) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (c *stdJsonCodec) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func (c *stdJsonCodec) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

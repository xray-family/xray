package codec

import (
	"encoding/json"
	"io"
)

var StdJsonCodec = new(stdJsonCodec)

type (
	Encoder interface {
		Encode(v interface{}) error
	}

	Decoder interface {
		Decode(v interface{}) error
	}

	Codec interface {
		NewEncoder(w io.Writer) Encoder
		NewDecoder(r io.Reader) Decoder
		Encode(v interface{}) ([]byte, error)
		Decode(data []byte, v interface{}) error
	}
)

type stdJsonCodec struct{}

func (c *stdJsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *stdJsonCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (c *stdJsonCodec) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func (c *stdJsonCodec) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

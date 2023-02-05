package uRouter

import (
	"encoding/json"
	"io"
)

const (
	ProtocolHTTP      = "http"
	ProtocolWebSocket = "websocket"
)

type (
	Any map[string]interface{}

	Encoder interface {
		Encode(v interface{}) error
	}

	Decoder interface {
		Decode(v interface{}) error
	}

	Codec interface {
		NewEncoder(w io.Writer) Encoder
		NewDecoder(r io.Reader) Decoder
	}
)

type stdJsonCodec struct{}

func (c *stdJsonCodec) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func (c *stdJsonCodec) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

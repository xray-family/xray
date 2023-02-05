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

	Form = MapHeader

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

func (c Any) Exists(k string) bool {
	_, ok := c[k]
	return ok
}

func (c Any) ToInt(k string) int {
	if v1, ok1 := c[k]; ok1 {
		if v2, ok2 := v1.(int); ok2 {
			return v2
		}
	}
	return 0
}

func (c Any) ToInt64(k string) int64 {
	if v1, ok1 := c[k]; ok1 {
		if v2, ok2 := v1.(int64); ok2 {
			return v2
		}
	}
	return 0
}

func (c Any) ToString(k string) string {
	if v1, ok1 := c[k]; ok1 {
		if v2, ok2 := v1.(string); ok2 {
			return v2
		}
	}
	return ""
}

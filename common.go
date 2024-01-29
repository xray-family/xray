package xray

import (
	"io"
)

type BytesReader interface {
	io.Reader
	Bytes() []byte
}

type (
	Any map[string]any

	Form map[string]string
)

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

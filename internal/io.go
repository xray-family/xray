package internal

import "io"

type (
	Closer interface {
		Close()
	}

	BytesReader interface {
		Bytes() []byte
	}

	Lener interface {
		Len() int
	}
)

func Close(source interface{}) error {
	if v, ok := source.(io.Closer); ok {
		return v.Close()
	}
	if v, ok := source.(Closer); ok {
		v.Close()
	}
	return nil
}

package uRouter

import (
	"bytes"
	"io"
)

type (
	Logger interface {
		Infof(format string, v ...interface{})
		Panicf(format string, v ...interface{})
	}

	BufferPool interface {
		SetSize(size int)
		Get() *bytes.Buffer
		Put(b *bytes.Buffer)
	}

	Header interface {
		Set(key, value string)
		Get(key string) string
		Del(key string)
		Len() int
		Range(f func(key, value string))
	}

	BytesReader interface {
		io.Reader
		Bytes() []byte
	}

	Closer interface {
		Close()
	}
)

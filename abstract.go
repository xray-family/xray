package uRouter

import (
	"bytes"
	"io"
)

type (
	LoggerInterface interface {
		Debug(format string, v ...interface{})
		Info(format string, v ...interface{})
		Warn(format string, v ...interface{})
		Error(format string, v ...interface{})
		Panic(format string, v ...interface{})
	}

	BufferPoolInterface interface {
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

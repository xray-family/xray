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

	BytesReader interface {
		io.Reader
		Bytes() []byte
	}

	Closer interface {
		Close()
	}
)

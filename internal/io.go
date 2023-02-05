package internal

import (
	"errors"
	"io"
)

type (
	Closer interface {
		Close()
	}

	BytesReader interface {
		io.Reader
		Bytes() []byte
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

func Write(w io.Writer, p []byte) error {
	num, err := w.Write(p)
	if err != nil {
		return err
	}
	if num != len(p) {
		return errors.New("io size error")
	}
	return nil
}

func Read(r io.Reader, p []byte) error {
	num, err := r.Read(p)
	if err != nil {
		return err
	}
	if num != len(p) {
		return errors.New("io size error")
	}
	return nil
}

func Copy(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}

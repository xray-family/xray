package internal

import "io"

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
	_, err := w.Write(p)
	return err
}

func Copy(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}

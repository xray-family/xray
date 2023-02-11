package internal

import (
	"github.com/lxzan/uRouter/constant"
	"io"
)

func Write(w io.Writer, p []byte) error {
	num, err := w.Write(p)
	if err != nil {
		return err
	}
	if num != len(p) {
		return constant.ErrIOBytesUnexpected
	}
	return nil
}

func Read(r io.Reader, p []byte) error {
	num, err := r.Read(p)
	if err != nil {
		return err
	}
	if num != len(p) {
		return constant.ErrIOBytesUnexpected
	}
	return nil
}

func Copy(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}

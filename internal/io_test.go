package internal

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type ioMocker struct{}

func (w ioMocker) Read(p []byte) (n int, err error) {
	num := len(p)
	if num <= 2 {
		return 1, nil
	}
	return 0, io.EOF
}

func (w ioMocker) Write(p []byte) (n int, err error) {
	num := len(p)
	if num <= 2 {
		return 1, nil
	}
	return 0, io.EOF
}

func TestWrite(t *testing.T) {
	var as = assert.New(t)

	t.Run("", func(t *testing.T) {
		var w = bytes.NewBufferString("")
		err := Write(w, []byte{})
		as.NoError(err)
	})

	t.Run("", func(t *testing.T) {
		var w = new(ioMocker)
		err := Write(w, []byte{'a', 'b'})
		as.Error(err)
	})

	t.Run("", func(t *testing.T) {
		var w = new(ioMocker)
		err := Write(w, []byte{'a', 'b', 'c'})
		as.Error(err)
	})
}

func TestRead(t *testing.T) {
	var as = assert.New(t)

	t.Run("", func(t *testing.T) {
		var r = bytes.NewBufferString("1234")
		err := Read(r, make([]byte, 2))
		as.NoError(err)
	})

	t.Run("", func(t *testing.T) {
		var r = new(ioMocker)
		err := Read(r, []byte{'a', 'b'})
		as.Error(err)
	})

	t.Run("", func(t *testing.T) {
		var r = new(ioMocker)
		err := Read(r, []byte{'a', 'b', 'c'})
		as.Error(err)
	})
}

func TestCopy(t *testing.T) {
	var a = bytes.NewBufferString("")
	var b = bytes.NewBufferString("12")
	assert.NoError(t, Copy(a, b))
}

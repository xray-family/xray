package internal

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type closerMocker struct{}

func (c *closerMocker) Close() {}

func TestClose(t *testing.T) {
	var as = assert.New(t)

	t.Run("", func(t *testing.T) {
		var c = io.NopCloser(bytes.NewBuffer(nil))
		as.NoError(Close(c))
	})

	t.Run("", func(t *testing.T) {
		var c = &closerMocker{}
		as.NoError(Close(c))
	})

	t.Run("", func(t *testing.T) {
		var c = bytes.NewBuffer(nil)
		as.NoError(Close(c))
	})
}

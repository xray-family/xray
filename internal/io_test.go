package internal

import (
	"bytes"
	"github.com/lxzan/gws"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestClose(t *testing.T) {
	var as = assert.New(t)

	t.Run("", func(t *testing.T) {
		var c = io.NopCloser(bytes.NewBuffer(nil))
		as.NoError(Close(c))
	})

	t.Run("", func(t *testing.T) {
		var c = &gws.Message{}
		as.NoError(Close(c))
	})

	t.Run("", func(t *testing.T) {
		var c = bytes.NewBuffer(nil)
		as.NoError(Close(c))
	})
}

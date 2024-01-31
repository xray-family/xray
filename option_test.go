package xray

import (
	"github.com/stretchr/testify/assert"
	"github.com/xray-family/xray/codec"
	"github.com/xray-family/xray/log"
	"testing"
	"time"
)

func TestWithGreeting(t *testing.T) {
	t.Run("", func(t *testing.T) {
		r := New()
		assert.True(t, r.conf.greeting.enabled)
		assert.Equal(t, r.conf.greeting.delay, time.Second)
	})

	t.Run("", func(t *testing.T) {
		r := New(WithGreeting(true, 2*time.Second))
		assert.True(t, r.conf.greeting.enabled)
		assert.Equal(t, r.conf.greeting.delay, 2*time.Second)
	})

	t.Run("", func(t *testing.T) {
		r := New(WithGreeting(false, 0))
		assert.False(t, r.conf.greeting.enabled)
	})
}

func TestWithJsonCodec(t *testing.T) {
	var conf = new(config)
	WithJsonCodec(codec.StdJsonCodec)(conf)
	assert.NotNil(t, conf.jsonCodec)
}

func TestWithLogger(t *testing.T) {
	var conf = new(config)
	WithLogger(log.StdLogger)(conf)
	assert.NotNil(t, conf.logger)
}

package xray

import (
	"github.com/xray-family/xray/codec"
	"github.com/xray-family/xray/log"
	"time"
)

type (
	config struct {
		greeting  greeting
		logger    log.Logger
		jsonCodec codec.Codec
	}

	greeting struct {
		enabled bool
		delay   time.Duration
	}
)

type Option func(c *config)

// WithGreeting 设置是否打印问候语, 已经打印的延迟时间
func WithGreeting(enabled bool, delay time.Duration) Option {
	return func(c *config) {
		c.greeting.enabled = enabled
		c.greeting.delay = delay
	}
}

func WithLogger(logger log.Logger) Option {
	return func(c *config) {
		c.logger = logger
	}
}

func WithJsonCodec(jsonCodec codec.Codec) Option {
	return func(c *config) {
		c.jsonCodec = jsonCodec
	}
}

func withInit() Option {
	return func(c *config) {
		if c.logger == nil {
			c.logger = log.StdLogger
		}
		if c.jsonCodec == nil {
			c.jsonCodec = codec.StdJsonCodec
		}
	}
}

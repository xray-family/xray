package xray

import "time"

type (
	config struct {
		greeting greeting
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

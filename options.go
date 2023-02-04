package uRouter

import (
	"github.com/go-playground/validator"
	"github.com/ugorji/go/codec"
)

var (
	defaultCodecHandle = new(codec.JsonHandle)
	defaultValidator   = validator.New()
)

type Config struct {
	CodecHandle codec.Handle
	Validator   *validator.Validate
}

type Option func(c *Config)

func withInitialize() Option {
	return func(c *Config) {
		if c.CodecHandle == nil {
			c.CodecHandle = defaultCodecHandle
		}
		if c.Validator == nil {
			c.Validator = defaultValidator
		}
	}
}

func WithCodecHandle(ch codec.Handle) Option {
	return func(c *Config) {
		c.CodecHandle = ch
	}
}

func WithValidator(v *validator.Validate) Option {
	return func(c *Config) {
		c.Validator = v
	}
}

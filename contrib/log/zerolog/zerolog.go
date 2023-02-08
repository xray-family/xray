package zerolog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var ZeroLogger = newLogger()

func newLogger() *Logger {
	return &Logger{
		caller: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}
}

type Logger struct {
	caller zerolog.Logger
}

func (c *Logger) Logger() zerolog.Logger {
	return c.caller
}

func (c *Logger) Debug(format string, v ...interface{}) {
	c.caller.Debug().Msgf(format, v...)
}

func (c *Logger) Info(format string, v ...interface{}) {
	c.caller.Info().Msgf(format, v...)
}

func (c *Logger) Warn(format string, v ...interface{}) {
	c.caller.Warn().Msgf(format, v...)
}

func (c *Logger) Error(format string, v ...interface{}) {
	c.caller.Error().Msgf(format, v...)
}

func (c *Logger) Panic(format string, v ...interface{}) {
	c.caller.Panic().Msgf(format, v...)
}

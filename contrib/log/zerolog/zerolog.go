package zerolog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var ZeroLogger = newLogger()

func newLogger() *zLogger {
	return &zLogger{
		caller: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}
}

type zLogger struct {
	caller zerolog.Logger
}

func (c *zLogger) UnWrap() zerolog.Logger {
	return c.caller
}

func (c *zLogger) Debug(format string, v ...interface{}) {
	c.caller.Debug().Msgf(format, v...)
}

func (c *zLogger) Info(format string, v ...interface{}) {
	c.caller.Info().Msgf(format, v...)
}

func (c *zLogger) Warn(format string, v ...interface{}) {
	c.caller.Warn().Msgf(format, v...)
}

func (c *zLogger) Error(format string, v ...interface{}) {
	c.caller.Error().Msgf(format, v...)
}

func (c *zLogger) Panic(format string, v ...interface{}) {
	c.caller.Panic().Msgf(format, v...)
}

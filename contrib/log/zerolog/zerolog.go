package zerolog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var ZeroLogger = newLogger()

func newLogger() *ZLogger { return &ZLogger{Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr})} }

type ZLogger struct {
	Logger zerolog.Logger
}

func (c *ZLogger) UnWrap() zerolog.Logger {
	return c.Logger
}

func (c *ZLogger) Debug(format string, v ...interface{}) {
	c.Logger.Debug().Msgf(format, v...)
}

func (c *ZLogger) Info(format string, v ...interface{}) {
	c.Logger.Info().Msgf(format, v...)
}

func (c *ZLogger) Warn(format string, v ...interface{}) {
	c.Logger.Warn().Msgf(format, v...)
}

func (c *ZLogger) Error(format string, v ...interface{}) {
	c.Logger.Error().Msgf(format, v...)
}

func (c *ZLogger) Panic(format string, v ...interface{}) {
	c.Logger.Panic().Msgf(format, v...)
}

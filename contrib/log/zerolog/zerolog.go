package zerolog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var Logger = newLogger()

func newLogger() *ZeroLogger {
	return &ZeroLogger{Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr})}
}

type ZeroLogger struct {
	Logger zerolog.Logger
}

func (c *ZeroLogger) UnWrap() zerolog.Logger {
	return c.Logger
}

func (c *ZeroLogger) Debug(format string, v ...any) {
	c.Logger.Debug().Msgf(format, v...)
}

func (c *ZeroLogger) Info(format string, v ...any) {
	c.Logger.Info().Msgf(format, v...)
}

func (c *ZeroLogger) Warn(format string, v ...any) {
	c.Logger.Warn().Msgf(format, v...)
}

func (c *ZeroLogger) Error(format string, v ...any) {
	c.Logger.Error().Msgf(format, v...)
}

func (c *ZeroLogger) Panic(format string, v ...any) {
	c.Logger.Panic().Msgf(format, v...)
}

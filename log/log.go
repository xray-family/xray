package log

import (
	"log"
)

type Logger interface {
	Debug(format string, v ...any)
	Info(format string, v ...any)
	Warn(format string, v ...any)
	Error(format string, v ...any)
	Panic(format string, v ...any)
}

var StdLogger = new(stdLogger)

type stdLogger struct{}

func (c *stdLogger) Debug(format string, v ...any) {
	log.Printf(format, v...)
}

func (c *stdLogger) Warn(format string, v ...any) {
	log.Printf(format, v...)
}

func (c *stdLogger) Error(format string, v ...any) {
	log.Printf(format, v...)
}

func (c *stdLogger) Info(format string, v ...any) {
	log.Printf(format, v...)
}

func (c *stdLogger) Panic(format string, v ...any) {
	log.Panicf(format, v...)
}

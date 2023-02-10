package uRouter

import (
	"log"
)

type logger struct{}

func (c *logger) Debug(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (c *logger) Warn(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (c *logger) Error(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (c *logger) Info(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (c *logger) Panic(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

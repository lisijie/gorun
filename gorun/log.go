package gorun

import (
	"fmt"
	"log"
)

var (
	NullLogger = &nullLogger{}
	StdLogger  = &stdLogger{}
)

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
}

type nullLogger struct{}

func (l *nullLogger) Debug(v ...interface{}) {}
func (l *nullLogger) Info(v ...interface{})  {}
func (l *nullLogger) Error(v ...interface{}) {}

type stdLogger struct{}

func (l *stdLogger) Debug(v ...interface{}) {
	log.Println("[debug]", fmt.Sprint(v...))
}

func (l *stdLogger) Info(v ...interface{}) {
	log.Println("[info]", fmt.Sprint(v...))
}
func (l *stdLogger) Error(v ...interface{}) {
	log.Println("[error]", fmt.Sprint(v...))
}

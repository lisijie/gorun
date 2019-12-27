package gorun

import (
	"fmt"
	"github.com/fatih/color"
	"log"
)

type Logger struct {
	isDebug bool
}

func (l *Logger) Debug(v ...interface{}) {
	if !l.isDebug {
		return
	}
	log.Println(color.GreenString("[DEBUG]"), fmt.Sprint(v...))
}

func (l *Logger) Info(v ...interface{}) {
	log.Println(color.BlueString("[INFO ]"), fmt.Sprint(v...))
}
func (l *Logger) Error(v ...interface{}) {
	log.Println(color.RedString("[ERROR]"), fmt.Sprint(v...))
}

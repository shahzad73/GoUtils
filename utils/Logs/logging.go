package Logs

import (
	"fmt"
	"time"
)

type Logger struct {
	timeFormat string
	debug      bool
}

func New(timeFormat string, debug bool) *Logger {
	return &Logger{
		timeFormat: timeFormat,
		debug:      debug,
	}
}

func (l *Logger) Log(s string) {
	if !l.debug {
		return
	}
	fmt.Printf("%s %s\n", time.Now().Format(l.timeFormat), s)
}

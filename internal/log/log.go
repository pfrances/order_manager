package log

import (
	"fmt"
	"io"
)

type Level int

const (
	Debug Level = iota
	Info
	Warning
	Error
)

type Logger struct {
	level Level
	w     io.Writer
	errW  io.Writer
}

func New(level Level, w io.Writer, errW io.Writer) *Logger {
	return &Logger{level: level, w: w, errW: errW}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.level <= Debug {
		l.w.Write([]byte("debug: " + fmt.Sprintf(format, args...)))
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.level <= Info {
		l.w.Write([]byte("info: " + fmt.Sprintf(format, args...)))
	}
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	if l.level <= Warning {
		l.errW.Write([]byte("warning: " + fmt.Sprintf(format, args...)))
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.level <= Error {
		l.errW.Write([]byte("error: " + fmt.Sprintf(format, args...)))
	}
}

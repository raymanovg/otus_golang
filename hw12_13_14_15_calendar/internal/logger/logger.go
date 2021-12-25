package logger

import (
	"fmt"
	"io"
)

type Logger struct {
	out   io.Writer
	level string
}

func New(out io.Writer, level string) *Logger {
	return &Logger{
		out:   out,
		level: level,
	}
}

func (l Logger) Info(msg string) {
	fmt.Fprintf(l.out, "[ INFO ] %s \n", msg)
}

func (l Logger) Warn(msg string) {
	fmt.Fprintf(l.out, "[ WARNINGN ] %s \n", msg)
}

func (l Logger) Error(msg string) {
	fmt.Fprintf(l.out, "[ ERROR ] %s \n", msg)
}

func (l Logger) Debug(msg string) {
	fmt.Fprintf(l.out, "[ DEBUG ] %s \n", msg)
}

func (l Logger) Write(msg string) {
	fmt.Fprintln(l.out, msg)
}

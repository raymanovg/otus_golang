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
	fmt.Fprintln(l.out, msg)
}

func (l Logger) Error(msg string) {
	// TODO
}

// TODO

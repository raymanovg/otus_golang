package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type ProgressBar struct {
	current, total int64
	percent        int
	rate           string
	graph          string
	pattern        string
	output         io.Writer
}

func NewProgressBar(total int64) *ProgressBar {
	return &ProgressBar{
		total:   total,
		current: 0,
		percent: 0,
		graph:   ">",
		pattern: "\r[%-100s]%3d%% %6d/%d bite",
		output:  os.Stdout,
	}
}

func (pb *ProgressBar) SetOutput(out io.Writer) {
	pb.output = out
}

func (pb *ProgressBar) SetTotal(t int64) {
	pb.total = t
}

func (pb *ProgressBar) SetCurrent(c int64) {
	pb.current = c
}

func (pb *ProgressBar) Start() {
	pb.Reset()
}

func (pb *ProgressBar) Finish() {
	fmt.Println()
}

func (pb *ProgressBar) Reset() {
	pb.current = 0
	pb.percent = 0
	pb.rate = ""
}

func (pb *ProgressBar) Add(n int64) {
	pb.current += n
	pb.write()
}

func (pb *ProgressBar) getPercent() int {
	return int(float32(pb.current) / float32(pb.total) * 100)
}

func (pb *ProgressBar) String() string {
	return fmt.Sprintf(pb.pattern, pb.rate, pb.percent, pb.total, pb.current)
}

func (pb *ProgressBar) write() {
	prev := pb.percent
	pb.percent = pb.getPercent()
	if pb.percent != prev && pb.percent%2 == 0 {
		pb.rate = strings.Repeat(pb.graph, pb.percent)
	}
	fmt.Fprint(pb.output, pb.String())
}

func (pb *ProgressBar) NewProgressBarWriter(writer io.Writer) *ProgressBarWriter {
	return &ProgressBarWriter{Writer: writer, bar: pb}
}

type ProgressBarWriter struct {
	io.Writer
	bar *ProgressBar
}

func (w *ProgressBarWriter) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)
	w.bar.Add(int64(n))
	return
}

func (w *ProgressBarWriter) Close() (err error) {
	if closer, ok := w.Writer.(io.Closer); ok {
		w.bar.Finish()
		return closer.Close()
	}
	return
}

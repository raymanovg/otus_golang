package main

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type ProgressBar struct {
	current, total int64
	percent        int
	rate           string
	graph          string
	pattern        string
}

func NewProgressBar(total int64) *ProgressBar {
	return &ProgressBar{
		total:   total,
		current: 0,
		percent: 0,
		graph:   ">",
		pattern: "\r[%-100s]%3d%% %8d/%d bite",
	}
}

func (pb *ProgressBar) SetTotal(t int64) *ProgressBar {
	pb.total = t
	return pb
}

func (pb *ProgressBar) SetCurrent(c int64) *ProgressBar {
	pb.current = c
	return pb
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

func (pb *ProgressBar) Add(n int) {
	pb.current += int64(n)
	pb.render()
}

func (pb *ProgressBar) getPercent() int {
	return int(float32(pb.current) / float32(pb.total) * 100)
}

func (pb *ProgressBar) render() {
	prev := pb.percent
	pb.percent = pb.getPercent()
	if pb.percent != prev && pb.percent%2 == 0 {
		pb.rate = strings.Repeat(pb.graph, pb.percent)
	}
	time.Sleep(time.Millisecond * 100)
	fmt.Printf(pb.pattern, pb.rate, pb.percent, pb.total, pb.current)
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
	w.bar.Add(n)
	return
}

func (w *ProgressBarWriter) Close() (err error) {
	if closer, ok := w.Writer.(io.Closer); ok {
		w.bar.Finish()
		return closer.Close()
	}
	return
}

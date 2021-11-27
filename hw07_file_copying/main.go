package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if from == "" {
		panic("src file is not defined")
	}
	if to == "" {
		panic("dst is not defined")
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		panic(fmt.Sprintf("unable to copy: %v", err))
	}
}

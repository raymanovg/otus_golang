package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("not enough args")
	}
	env, err := ReadDir(args[0])
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(RunCmd(args[1:], env))
}

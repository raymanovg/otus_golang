package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout string

	logger = log.New(os.Stderr, "", 0)
)

func init() {
	flag.StringVar(&timeout, "timeout", "10s", "timeout for connection")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		logger.Println("undefined host and port")
		return
	}

	timeoutDur, err := time.ParseDuration(timeout)
	if err != nil {
		logger.Println("unable to parse duration")
		return
	}

	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, timeoutDur, os.Stdin, os.Stdout)
	if err = client.Connect(); err != nil {
		logger.Printf("unable to connect to server %s\n", address)
		return
	}

	logger.Printf("...connected to %s\n", address)

	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	go send(cancel, client)
	go receive(cancel, client)

	<-ctx.Done()
}

func send(cancel context.CancelFunc, client TelnetClient) {
	defer cancel()
	if err := client.Send(); err != nil {
		logger.Printf("unexpected sending err: %v", err)
		return
	}
	logger.Println("...EOF")
}

func receive(cancel context.CancelFunc, client TelnetClient) {
	defer cancel()
	if err := client.Receive(); err != nil {
		logger.Printf("unexpected receiving error: %v", err)
		return
	}
	logger.Println("...connection was closed by peer")
}

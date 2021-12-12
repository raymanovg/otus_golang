package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout string

func init() {
	flag.StringVar(&timeout, "timeout", "10s", "timeout for connection")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "undefined host and port")
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)

	timeoutDur, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse duration")
		return
	}

	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, timeoutDur, os.Stdin, os.Stdout)
	connErr := client.Connect()
	if connErr != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to server %s\n", address)
		return
	}

	fmt.Fprintf(os.Stderr, "...connected to %s\n", address)

	defer client.Close()

	go send(ctx, cancel, client)
	go receive(ctx, cancel, client)

	<-ctx.Done()
}

func send(ctx context.Context, cancel context.CancelFunc, client TelnetClient) {
	defer cancel()
	for {
		err := client.Send()
		select {
		case <-ctx.Done():
			return
		default:
			if err == nil {
				continue
			}
			if err == io.EOF {
				fmt.Fprintln(os.Stderr, "...EOF")
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}
	}
}

func receive(ctx context.Context, cancel context.CancelFunc, client TelnetClient) {
	defer cancel()
	for {
		err := client.Receive()
		select {
		case <-ctx.Done():
			return
		default:
			if err == nil {
				continue
			}
			if err == ErrConnectionClosedByPeer {
				fmt.Fprintln(os.Stderr, "...connection was closed by peer")
			}
			return
		}
	}
}

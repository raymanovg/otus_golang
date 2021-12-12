package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	ErrConnectionClosedByPeer = errors.New("connection was closed by peer")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Telnet struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *Telnet) Close() error {
	t.in.Close()
	return t.conn.Close()
}

func (t *Telnet) Connect() error {
	dialer := net.Dialer{
		Timeout: t.timeout,
	}
	ctx := context.Background()
	conn, err := dialer.DialContext(ctx, "tcp", t.address)
	if err != nil {
		return fmt.Errorf("unble to connect: %w", err)
	}
	t.conn = conn
	return nil
}

func (t *Telnet) Send() error {
	scanner := bufio.NewScanner(t.in)
	if !scanner.Scan() {
		return io.EOF
	}
	in := append(scanner.Bytes(), []byte("\n")...)
	_, err := t.conn.Write(in)
	if err != nil {
		return fmt.Errorf("unable to send: %w", err)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("unable to read: %w", err)
	}
	return nil
}

func (t *Telnet) Receive() error {
	scanner := bufio.NewScanner(t.conn)
	if !scanner.Scan() {
		return ErrConnectionClosedByPeer
	}
	out := append(scanner.Bytes(), []byte("\n")...)
	_, err := t.out.Write(out)
	if err != nil {
		return fmt.Errorf("unable to write: %w", err)
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

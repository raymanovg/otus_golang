package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrConnectionClosedByPeer = errors.New("connection was closed by peer")

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
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
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
	if _, err := t.conn.Write(in); err != nil {
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
	if _, err := t.out.Write(out); err != nil {
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

package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrConnectionNotEstablished = errors.New("connection not established")

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
	if t.conn == nil {
		return ErrConnectionNotEstablished
	}
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return fmt.Errorf("unable to send: %w", err)
	}
	return nil
}

func (t *Telnet) Receive() error {
	if t.conn == nil {
		return ErrConnectionNotEstablished
	}
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return fmt.Errorf("unable to receive: %w", err)
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

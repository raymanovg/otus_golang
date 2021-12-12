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
	return copyBytes(t.in, t.conn)
}

func (t *Telnet) Receive() error {
	if t.conn == nil {
		return ErrConnectionNotEstablished
	}
	return copyBytes(t.conn, t.out)
}

func copyBytes(src io.Reader, dest io.Writer) error {
	if _, err := io.Copy(dest, src); err != nil {
		return fmt.Errorf("unable to copyBytes: %w", err)
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

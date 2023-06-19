package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	ErrErrorsConnect = errors.New("unable to connect")
	ErrErrorsSend    = errors.New("unable to send")
	ErrErrorsReceive = errors.New("unable to receive")
	ErrErrorsClose   = errors.New("error with close")
	ErrErrorsEOF     = errors.New("eof")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Telnet struct {
	in          io.ReadCloser
	out         io.Writer
	timeout     time.Duration
	address     string
	conn        net.Conn
	isconnected bool
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		in:      in,
		out:     out,
		timeout: timeout,
		address: address,
	}
}

func (t *Telnet) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), ErrErrorsConnect)
	}
	t.conn = conn
	t.isconnected = true
	return nil
}

func (t *Telnet) Send() error {
	_, err := io.Copy(t.conn, t.in)
	if err != nil && t.isconnected {
		return fmt.Errorf("%s %w", err.Error(), ErrErrorsSend)
	}
	return nil
}

func (t *Telnet) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	if err != nil && t.isconnected {
		return fmt.Errorf("%s %w", err.Error(), ErrErrorsReceive)
	}
	return nil
}

func (t *Telnet) Close() error {
	if !t.isconnected {
		return nil
	}
	t.isconnected = false
	err := t.conn.Close()
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), ErrErrorsClose)
	}
	return nil
}

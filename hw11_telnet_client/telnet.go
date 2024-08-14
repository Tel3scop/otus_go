package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type SimpleTelnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &SimpleTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (tc *SimpleTelnetClient) Connect() error {
	var d net.Dialer
	var err error
	ctx := context.Background()
	tc.conn, err = d.DialContext(ctx, "tcp", tc.address)
	if err != nil {
		return err
	}
	return nil
}

func (tc *SimpleTelnetClient) Close() error {
	if tc.conn != nil {
		return tc.conn.Close()
	}
	return nil
}

func (tc *SimpleTelnetClient) Send() error {
	if tc.conn == nil {
		return fmt.Errorf("not connected")
	}
	_, err := io.Copy(tc.conn, tc.in)
	return err
}

func (tc *SimpleTelnetClient) Receive() error {
	if tc.conn == nil {
		return fmt.Errorf("not connected")
	}
	_, err := io.Copy(tc.out, tc.conn)
	return err
}

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
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
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", tc.address)
	return nil
}

func (tc *SimpleTelnetClient) Close() error {
	if tc.conn != nil {
		err := tc.conn.Close()
		fmt.Fprintln(os.Stderr, "...Connection closed")
		return err
	}
	return nil
}

func (tc *SimpleTelnetClient) Send() error {
	if tc.conn == nil {
		return fmt.Errorf("not connected")
	}
	_, err := io.Copy(tc.conn, tc.in)
	if errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, "...EOF")
		return err
	}
	return err
}

func (tc *SimpleTelnetClient) Receive() error {
	if tc.conn == nil {
		return fmt.Errorf("not connected")
	}
	_, err := io.Copy(tc.out, tc.conn)
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		fmt.Fprintln(os.Stderr, "...Connection timed out")
	}
	return err
}

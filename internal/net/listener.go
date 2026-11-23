package net

import (
	"fmt"
	"net"
)

// TCPListener accepts connections and wraps them as transports.
type TCPListener struct {
	ln net.Listener
}

// Listen starts listening on addr, e.g. ":7777".
func Listen(addr string) (*TCPListener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("could not listen on %s: %w", addr, err)
	}
	return &TCPListener{ln: ln}, nil
}

func (l *TCPListener) Accept() (Transport, error) {
	conn, err := l.ln.Accept()
	if err != nil {
		return nil, err
	}
	return NewTransport(conn), nil
}

func (l *TCPListener) Addr() string {
	return l.ln.Addr().String()
}

func (l *TCPListener) Close() error {
	return l.ln.Close()
}

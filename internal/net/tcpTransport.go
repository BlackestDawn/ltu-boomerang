package net

import (
	"bufio"
	"encoding/json"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

// DialTimeout is how long Dial waits for the host to answer.
const DialTimeout = 10 * time.Second

// WriteTimeout is how long Send waits for the other side to take a message,
// so a client that stopped reading can not stall the game.
const WriteTimeout = 10 * time.Second

// TCPTransport sends messages as newline separated JSON over a connection.
type TCPTransport struct {
	conn net.Conn

	sendMu sync.Mutex // serialises writes, Send is called from several goroutines
	enc    *json.Encoder
	dec    *json.Decoder
}

// NewTransport wraps an established connection. Any net.Conn works, which
// lets tests use net.Pipe.
func NewTransport(conn net.Conn) *TCPTransport {
	return &TCPTransport{
		conn: conn,
		enc:  json.NewEncoder(conn), // Encode ends every value with a newline
		dec:  json.NewDecoder(bufio.NewReader(conn)),
	}
}

// Dial connects to a host.
func Dial(remote string) (*TCPTransport, error) {
	if !strings.Contains(remote, ":") {
		remote += ":" + core.ListenPort
	}
	conn, err := net.DialTimeout("tcp", remote, DialTimeout)
	if err != nil {
		return nil, err
	}
	return NewTransport(conn), nil
}

func (t *TCPTransport) Send(m Message) error {
	t.sendMu.Lock()
	defer t.sendMu.Unlock()
	if err := t.conn.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
		return err
	}
	return t.enc.Encode(m)
}

func (t *TCPTransport) Receive() (Message, error) {
	var m Message
	err := t.dec.Decode(&m)
	return m, err
}

// SetReadDeadline makes Receive fail if nothing arrives before t. The zero
// time removes the deadline.
func (t *TCPTransport) SetReadDeadline(deadline time.Time) error {
	return t.conn.SetReadDeadline(deadline)
}

func (t *TCPTransport) Close() error {
	return t.conn.Close()
}

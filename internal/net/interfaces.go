package net

import "time"

// Transport is one connection to the other side, sending and receiving
// whole messages. Send may be called from several goroutines at once, Receive
// from one at a time.
type Transport interface {
	Send(m Message) error
	Receive() (Message, error)
	// SetReadDeadline makes Receive fail if nothing arrives before the
	// deadline. The zero time removes it.
	SetReadDeadline(deadline time.Time) error
	Close() error
}

// Listener accepts incoming connections on the host.
type Listener interface {
	Accept() (Transport, error)
	Addr() string
	Close() error
}

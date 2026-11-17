package net

type Transport interface {
	Send(m Message) error
	Recieve() (Message, error)
}

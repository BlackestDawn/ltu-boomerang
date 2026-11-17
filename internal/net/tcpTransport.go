package net

type TCPTransport struct {
	remote string
}

func (t *TCPTransport) Send(m Message) error {
	return nil
}

func (t *TCPTransport) Recieve() (Message, error) {
	return Message{}, nil
}

package core

import "context"

// subscriberBuffer is how many events a subscriber can fall behind before the
// server (and therefore the event source) blocks.
const subscriberBuffer = 32

type EventServer interface {
	Subscribe() <-chan Event
	CancelSubscription(<-chan Event)
}

type eventServer struct {
	ctx            context.Context
	source         <-chan Event
	listeners      []chan Event
	addListener    chan chan Event
	removeListener chan (<-chan Event)
}

// Subscribe returns a channel receiving every event from the source. The
// channel is closed when the server stops (context cancelled or source closed).
func (e *eventServer) Subscribe() <-chan Event {
	newListener := make(chan Event, subscriberBuffer)
	select {
	case e.addListener <- newListener:
	case <-e.ctx.Done():
		close(newListener)
	}
	return newListener
}

func (e *eventServer) CancelSubscription(channel <-chan Event) {
	select {
	case e.removeListener <- channel:
	case <-e.ctx.Done():
	}
}

func NewEventServer(ctx context.Context, source <-chan Event) EventServer {
	service := &eventServer{
		ctx:            ctx,
		source:         source,
		listeners:      make([]chan Event, 0, 16),
		addListener:    make(chan chan Event),
		removeListener: make(chan (<-chan Event)),
	}
	go service.serve()
	return service
}

func (e *eventServer) serve() {
	defer func() {
		for _, listener := range e.listeners {
			close(listener)
		}
	}()

	for {
		select {
		case <-e.ctx.Done():
			return
		case newListener := <-e.addListener:
			e.listeners = append(e.listeners, newListener)
		case listenerToRemove := <-e.removeListener:
			e.remove(listenerToRemove)
		case val, ok := <-e.source:
			if !ok {
				return
			}
			for _, listener := range e.listeners {
				select {
				case listener <- val:
				case <-e.ctx.Done():
					return
				}
			}
		}
	}
}

func (e *eventServer) remove(target <-chan Event) {
	for i, ch := range e.listeners {
		if ch == target {
			e.listeners[i] = e.listeners[len(e.listeners)-1]
			e.listeners = e.listeners[:len(e.listeners)-1]
			close(ch)
			return
		}
	}
}

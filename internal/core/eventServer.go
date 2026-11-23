package core

import "context"

// subscriberBuffer is how many events a subscriber can fall behind before the
// server (and therefore the publisher) blocks.
const subscriberBuffer = 32

// Publisher is the sending side of the event server.
type Publisher interface {
	// Publish returns once ev is queued on every subscriber channel. Events
	// are queued in the order they are published, so anything a goroutine
	// does after Publish returns happens after the event is in every queue.
	Publish(ev Event)
}

type EventServer interface {
	Publisher
	Subscribe() <-chan Event
	CancelSubscription(<-chan Event)
}

type publishRequest struct {
	ev   Event
	done chan struct{}
}

type eventServer struct {
	ctx            context.Context
	publish        chan publishRequest
	listeners      []chan Event
	addListener    chan chan Event
	removeListener chan (<-chan Event)
}

// Subscribe returns a channel receiving every event published from now on.
// The channel is closed when the server stops (context cancelled) or the
// subscription is cancelled.
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

func (e *eventServer) Publish(ev Event) {
	req := publishRequest{ev: ev, done: make(chan struct{})}
	select {
	case e.publish <- req:
	case <-e.ctx.Done():
		return
	}
	select {
	case <-req.done:
	case <-e.ctx.Done():
	}
}

// NewEventServer starts a server that runs until ctx is cancelled.
func NewEventServer(ctx context.Context) EventServer {
	service := &eventServer{
		ctx:            ctx,
		publish:        make(chan publishRequest),
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
		case req := <-e.publish:
			for _, listener := range e.listeners {
				select {
				case listener <- req.ev:
				case <-e.ctx.Done():
					return
				}
			}
			close(req.done)
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

package agent

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

// eventPump hands events from a subscription to handle, one at a time on its
// own goroutine, until the subscription is closed.
//
// Controllers call Flush before asking their player anything. Together with
// core.Publisher's guarantee that published events are already queued, this
// makes sure the player has seen every event that happened before the
// question, e.g. the round scores before the next round's throw.
type eventPump struct {
	events  <-chan core.Event
	handle  func(core.Event)
	flushes chan chan struct{}
	done    chan struct{}
}

func newEventPump(events <-chan core.Event, handle func(core.Event)) *eventPump {
	p := &eventPump{
		events:  events,
		handle:  handle,
		flushes: make(chan chan struct{}),
		done:    make(chan struct{}),
	}
	go p.run()
	return p
}

func (p *eventPump) run() {
	defer close(p.done)
	for {
		select {
		case e, ok := <-p.events:
			if !ok {
				return
			}
			p.handle(e)
		case ack := <-p.flushes:
			open := p.drain()
			close(ack)
			if !open {
				return
			}
		}
	}
}

// drain handles every event already queued. It reports false if the
// subscription was closed.
func (p *eventPump) drain() bool {
	for {
		select {
		case e, ok := <-p.events:
			if !ok {
				return false
			}
			p.handle(e)
		default:
			return true
		}
	}
}

// Flush returns once every event queued before the call has been handled.
func (p *eventPump) Flush() {
	ack := make(chan struct{})
	select {
	case p.flushes <- ack:
		<-ack
	case <-p.done:
	}
}

// Wait returns once the subscription is closed and every event handled.
func (p *eventPump) Wait() {
	<-p.done
}

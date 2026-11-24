package core

import (
	"context"
	"testing"
	"time"
)

// Every subscriber gets every event, in publish order.
func TestEventServerDeliversInOrder(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	es := NewEventServer(ctx)
	a, b := es.Subscribe(), es.Subscribe()

	for i := range 5 {
		es.Publish(RoundBeginsEvent(i))
	}
	for _, ch := range []<-chan Event{a, b} {
		for i := range 5 {
			select {
			case ev := <-ch:
				if ev != RoundBeginsEvent(i) {
					t.Fatalf("got %v, want round %d", ev, i)
				}
			case <-time.After(time.Second):
				t.Fatal("timed out waiting for event")
			}
		}
	}
}

func TestEventServerCancelClosesChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	es := NewEventServer(ctx)
	ch := es.Subscribe()
	es.CancelSubscription(ch)

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("received event after cancelling")
		}
	case <-time.After(time.Second):
		t.Fatal("channel not closed after cancelling")
	}
}

func TestMarshalUnknownEventFails(t *testing.T) {
	if _, err := UnmarshalEvent([]byte(`{"type":"bogus","data":{}}`)); err == nil {
		t.Error("expected error for unknown event type")
	}
}

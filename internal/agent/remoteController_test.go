package agent

import (
	"context"
	"math/rand"
	gonet "net"
	"slices"
	"testing"
	"time"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/net"
	"github.com/BlackestDawn/ltu-boomerang/internal/scoring"
)

func newTestRemote(t *testing.T) (*RemoteController, *net.TCPTransport, core.EventServer) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	es := core.NewEventServer(ctx)

	hostSide, clientSide := gonet.Pipe()
	client := net.NewTransport(clientSide)
	t.Cleanup(func() { client.Close() })
	c := NewRemoteController(net.NewTransport(hostSide), es, NewBotController(rand.New(rand.NewSource(1))))
	return c, client, es
}

// receive reads messages until one of type want arrives, failing on timeout.
func receive(t *testing.T, client *net.TCPTransport, want net.MessageType) net.Message {
	t.Helper()
	client.SetReadDeadline(time.Now().Add(2 * time.Second))
	for {
		msg, err := client.Receive()
		if err != nil {
			t.Fatalf("waiting for %s: %v", want, err)
		}
		if msg.Type == want {
			return msg
		}
	}
}

func answer(t *testing.T, client *net.TCPTransport, choice int) {
	t.Helper()
	msg, _ := net.NewMessage(net.Answer, net.AnswerData{Choice: choice})
	if err := client.Send(msg); err != nil {
		t.Fatalf("send answer: %v", err)
	}
}

func testHand() []*core.Card {
	swim := &core.Symbol{Name: "Swimming"}
	walk := &core.Symbol{Name: "Bushwalking"}
	return []*core.Card{
		{Name: "A", Symbols: map[core.Category]*core.Symbol{core.Activity: swim}},
		{Name: "B", Symbols: map[core.Category]*core.Symbol{core.Activity: walk}},
		{Name: "C", Symbols: map[core.Category]*core.Symbol{core.Activity: swim}},
	}
}

func TestRemoteChooseDraftReturnsHostCard(t *testing.T) {
	c, client, _ := newTestRemote(t)
	hand := testHand()

	got := make(chan *core.Card)
	go func() { got <- c.ChooseDraft(hand) }()

	var ask net.AskCardData
	receive(t, client, net.AskDraw).Decode(&ask)
	if len(ask.Hand) != 3 || ask.Hand[2].Name != "C" {
		t.Fatalf("unexpected hand sent: %+v", ask.Hand)
	}

	// Out of range is refused and the question repeated.
	answer(t, client, 7)
	receive(t, client, net.Error)
	receive(t, client, net.AskDraw)
	answer(t, client, 2)

	if card := <-got; card != hand[2] {
		t.Fatalf("got %v, want the host's own pointer to card C", card)
	}
}

func TestRemoteChooseActivity(t *testing.T) {
	c, client, _ := newTestRemote(t)
	hand := testHand()
	swim := hand[0].Symbols[core.Activity]

	got := make(chan *core.Symbol)
	go func() { got <- c.ChooseActivity(hand, nil) }()

	var ask net.AskActivityData
	receive(t, client, net.AskActivity).Decode(&ask)
	// Options are sorted by name: Bushwalking, Swimming.
	if len(ask.Options) != 2 || ask.Options[1].Name != "Swimming" || ask.Counts[1] != 2 {
		t.Fatalf("unexpected options: %+v %v", ask.Options, ask.Counts)
	}
	answer(t, client, 1)
	if s := <-got; s != swim {
		t.Fatalf("got %v, want swimming", s)
	}

	go func() { got <- c.ChooseActivity(hand, nil) }()
	receive(t, client, net.AskActivity)
	answer(t, client, -1)
	if s := <-got; s != nil {
		t.Fatalf("skip gave %v, want nil", s)
	}
}

func TestRemoteFallsBackToBotWhenLost(t *testing.T) {
	c, client, es := newTestRemote(t)
	client.Close()

	hand := testHand()
	done := make(chan *core.Card)
	go func() { done <- c.ChooseThrow(hand) }()
	select {
	case card := <-done:
		if !slices.Contains(hand, card) {
			t.Fatalf("fallback picked %v, not from hand", card)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ChooseThrow blocked on a closed connection")
	}

	// Events must still be drained so the event server never blocks.
	published := make(chan struct{})
	go func() {
		for i := range 64 {
			es.Publish(core.RoundBeginsEvent(i))
		}
		close(published)
	}()
	select {
	case <-published:
	case <-time.After(2 * time.Second):
		t.Fatal("event server blocked after connection loss")
	}
	c.Finish(scoring.SessionResult{})
}

func TestRemoteFinishSendsQueuedEventsThenFinal(t *testing.T) {
	c, client, es := newTestRemote(t)
	es.Publish(core.RoundBeginsEvent(0))

	msg := receive(t, client, net.Event)
	if ev, err := core.UnmarshalEvent(msg.Data); err != nil || ev != core.RoundBeginsEvent(0) {
		t.Fatalf("got event %v (%v)", ev, err)
	}

	sheet := &core.ScoreSheet{CategoryScores: map[core.Category][]int{core.Site: {5}}}
	go c.Finish(scoring.SessionResult{
		Scores: []*core.ScoreSheet{sheet},
		Winner: 0,
	})

	var final net.FinalData
	receive(t, client, net.Final).Decode(&final)
	if len(final.ScoreSheets) != 1 || final.ScoreSheets[0].Total() != 5 {
		t.Fatalf("unexpected final: %+v", final)
	}
	receive(t, client, net.Bye)
}

// Events published before a question must reach the client before it, or the
// client is asked for its next throw before seeing the round scores.
func TestRemoteSendsEarlierEventsBeforeAsk(t *testing.T) {
	for range 50 {
		c, client, es := newTestRemote(t)
		hand := testHand()

		es.Publish(core.RoundScoreEvent{0: {core.Site: 3}})
		es.Publish(core.RoundBeginsEvent(1))
		go c.ChooseThrow(hand)

		client.SetReadDeadline(time.Now().Add(2 * time.Second))
		var order []net.MessageType
		for len(order) < 3 {
			msg, err := client.Receive()
			if err != nil {
				t.Fatalf("receive: %v", err)
			}
			order = append(order, msg.Type)
		}
		want := []net.MessageType{net.Event, net.Event, net.AskThrow}
		if !slices.Equal(order, want) {
			t.Fatalf("got messages %v, want %v", order, want)
		}
		answer(t, client, 0)
	}
}

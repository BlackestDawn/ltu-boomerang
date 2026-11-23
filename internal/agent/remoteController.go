package agent

import (
	"errors"
	"fmt"
	"sync"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/net"
	"github.com/BlackestDawn/ltu-boomerang/internal/scoring"
)

// errLost is returned by ask once the connection is gone.
var errLost = errors.New("connection to remote player lost")

// RemoteController plays a seat for a client connected over the network. It
// forwards every game event to the client and asks it for each choice. If the
// connection is lost the fallback bot takes over the seat for the rest of the
// game.
type RemoteController struct {
	net          net.Transport
	es           core.EventServer
	subscription <-chan core.Event
	events       *eventPump
	fallback     *BotController

	mu   sync.Mutex // guards lost
	lost bool
}

// NewRemoteController subscribes to es and starts forwarding events to the
// client on transport. Call Finish when the game is over.
func NewRemoteController(transport net.Transport, es core.EventServer, fallback *BotController) *RemoteController {
	c := &RemoteController{
		net:          transport,
		es:           es,
		subscription: es.Subscribe(),
		fallback:     fallback,
	}
	c.events = newEventPump(c.subscription, c.forwardEvent)
	return c
}

// forwardEvent sends an event to the client. Once the connection is lost
// events are dropped, but still taken off the subscription so the event
// server never blocks on it.
func (c *RemoteController) forwardEvent(e core.Event) {
	if c.isLost() {
		return
	}
	msg, err := net.NewEventMessage(e)
	if err != nil {
		c.sendError(err)
		return
	}
	if err := c.net.Send(msg); err != nil {
		c.markLost()
	}
}

func (c *RemoteController) ChooseThrow(hand []*core.Card) *core.Card {
	choice, err := c.ask(net.AskThrow, net.AskCardData{Hand: core.DerefSlice(hand)}, 0, len(hand)-1)
	if err != nil {
		return c.fallback.ChooseThrow(hand)
	}
	return hand[choice]
}

func (c *RemoteController) ChooseDraft(hand []*core.Card) *core.Card {
	choice, err := c.ask(net.AskDraw, net.AskCardData{Hand: core.DerefSlice(hand)}, 0, len(hand)-1)
	if err != nil {
		return c.fallback.ChooseDraft(hand)
	}
	return hand[choice]
}

func (c *RemoteController) ChooseActivity(hand []*core.Card, alreadyScored []*core.Symbol) *core.Symbol {
	options, counts := activityOptions(hand, alreadyScored)
	if len(options) == 0 {
		return nil
	}
	payload := net.AskActivityData{
		Draft:   core.DerefSlice(hand),
		Options: core.DerefSlice(options),
		Counts:  counts,
	}
	choice, err := c.ask(net.AskActivity, payload, -1, len(options)-1)
	if err != nil {
		return c.fallback.ChooseActivity(hand, alreadyScored)
	}
	if choice < 0 {
		return nil
	}
	return options[choice]
}

// Finish stops forwarding events, sends the final results and closes the
// connection. Events still queued for the client are sent first.
func (c *RemoteController) Finish(result scoring.SessionResult) {
	c.es.CancelSubscription(c.subscription)
	c.events.Wait()

	if !c.isLost() {
		final := net.FinalData{Winner: result.Winner}
		for _, s := range result.Ranking {
			final.ScoreSheets = append(final.ScoreSheets, *s.ScoreSheet)
		}
		if msg, err := net.NewMessage(net.Final, final); err == nil {
			c.net.Send(msg)
		}
		c.net.Send(net.Message{Type: net.Bye})
	}
	c.net.Close()
}

// ask sends a question and waits for an answer in [min, max]. Events
// published before the call are sent first. Answers out of range are refused
// and the client is asked again. Any other message is ignored.
func (c *RemoteController) ask(t net.MessageType, payload any, min, max int) (int, error) {
	c.events.Flush()
	if c.isLost() {
		return 0, errLost
	}
	msg, err := net.NewMessage(t, payload)
	if err != nil {
		return 0, err
	}
	if err := c.net.Send(msg); err != nil {
		c.markLost()
		return 0, errLost
	}

	for {
		reply, err := c.net.Receive()
		if err != nil {
			c.markLost()
			return 0, errLost
		}
		switch reply.Type {
		case net.Answer:
			var a net.AnswerData
			if err := reply.Decode(&a); err != nil {
				c.sendError(fmt.Errorf("could not read answer: %w", err))
			} else if a.Choice < min || a.Choice > max {
				c.sendError(fmt.Errorf("answer %d out of range %d-%d", a.Choice, min, max))
			} else {
				return a.Choice, nil
			}
			// Repeat the question so the client knows it is still expected.
			if err := c.net.Send(msg); err != nil {
				c.markLost()
				return 0, errLost
			}
		case net.Bye:
			c.markLost()
			return 0, errLost
		}
	}
}

func (c *RemoteController) sendError(err error) {
	msg, mErr := net.NewMessage(net.Error, net.ErrorData{Message: err.Error()})
	if mErr != nil {
		return
	}
	if c.net.Send(msg) != nil {
		c.markLost()
	}
}

// markLost switches the seat over to the fallback bot and closes the
// connection, which also unblocks a pending Receive.
func (c *RemoteController) markLost() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lost {
		return
	}
	c.lost = true
	c.net.Close()
}

func (c *RemoteController) isLost() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lost
}

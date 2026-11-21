package agent

import (
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/net"
)

type RemoteController struct {
	net      net.Transport
	eventBus <-chan core.Event
}

func NewRemoteController(ev <-chan core.Event) *RemoteController {
	return &RemoteController{
		net:      &net.TCPTransport{},
		eventBus: ev,
	}
}

func (c *RemoteController) ChooseThrow(hand []*core.Card) *core.Card {
	return hand[0]
}

func (c *RemoteController) ChooseDraft(hand []*core.Card) *core.Card {
	return hand[0]
}

func (c *RemoteController) ChooseActivity(hand []*core.Card, alreadyScored []*core.Symbol) *core.Symbol {
	return nil
}

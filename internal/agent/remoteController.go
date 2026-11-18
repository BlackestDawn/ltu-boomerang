package agent

import (
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/net"
)

type RemoteController struct {
	net net.Transport
}

func NewRemoteController() *RemoteController {
	return &RemoteController{
		net: &net.TCPTransport{},
	}
}

func (c *RemoteController) ChooseThrow(hand []*core.Card) *core.Card {
	return &core.Card{}
}

func (c *RemoteController) ChooseDraft(hand []*core.Card) *core.Card {
	return &core.Card{}
}

func (c *RemoteController) ChooseActivity(hand []*core.Card, alreadyScored []*core.Symbol) *core.Symbol {
	return nil
}

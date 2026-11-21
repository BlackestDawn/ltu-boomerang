package agent

import (
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/ui"
)

type HumanController struct {
	ui       *ui.TextUI
	eventBus <-chan core.Event
}

// NewHumanController creates a controller and starts forwarding events from
// ev to the UI. Forwarding stops when ev is closed.
func NewHumanController(ui *ui.TextUI, ev <-chan core.Event) *HumanController {
	c := &HumanController{
		ui:       ui,
		eventBus: ev,
	}
	go c.forwardEvents()
	return c
}

func (c *HumanController) forwardEvents() {
	for e := range c.eventBus {
		c.ui.OnEvent(e)
	}
}

func (c *HumanController) ChooseThrow(hand []*core.Card) *core.Card {
	num := c.ui.AskForNumber(len(hand), 1)
	return hand[num-1]
}

func (c *HumanController) ChooseDraft(hand []*core.Card) *core.Card {
	num := c.ui.AskForNumber(len(hand), 1)
	return hand[num-1]
}

func (c *HumanController) ChooseActivity(hand []*core.Card, alreadyScored []*core.Symbol) *core.Symbol {
	return nil
}

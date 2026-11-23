package agent

import (
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/ui"
)

type HumanController struct {
	ui     *ui.TextUI
	events *eventPump
}

// NewHumanController creates a controller and starts forwarding events from
// ev to the UI. Forwarding stops when ev is closed.
func NewHumanController(ui *ui.TextUI, ev <-chan core.Event) *HumanController {
	return &HumanController{
		ui:     ui,
		events: newEventPump(ev, ui.OnEvent),
	}
}

func (c *HumanController) ChooseThrow(hand []*core.Card) *core.Card {
	c.events.Flush()
	return hand[c.ui.AskForCard(core.DerefSlice(hand), "Choose your throw card")]
}

func (c *HumanController) ChooseDraft(hand []*core.Card) *core.Card {
	c.events.Flush()
	return hand[c.ui.AskForCard(core.DerefSlice(hand), "Choose a card to draft")]
}

func (c *HumanController) ChooseActivity(hand []*core.Card, alreadyScored []*core.Symbol) *core.Symbol {
	options, counts := activityOptions(hand, alreadyScored)
	if len(options) == 0 {
		return nil
	}
	c.events.Flush()
	choice := c.ui.AskForActivity(core.DerefSlice(hand), core.DerefSlice(options), counts)
	if choice < 0 {
		return nil
	}
	return options[choice]
}

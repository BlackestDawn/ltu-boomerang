package agent

import (
	"slices"
	"strings"

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
	return hand[c.ui.AskForCard(derefCards(hand), "Choose your throw card")]
}

func (c *HumanController) ChooseDraft(hand []*core.Card) *core.Card {
	return hand[c.ui.AskForCard(derefCards(hand), "Choose a card to draft")]
}

func (c *HumanController) ChooseActivity(hand []*core.Card, alreadyScored []*core.Symbol) *core.Symbol {
	var options []*core.Symbol
	counts := core.FilterSymbolCategoryWithCount(hand, core.Activity)
	for symbol := range counts {
		if !slices.Contains(alreadyScored, symbol) {
			options = append(options, symbol)
		}
	}
	if len(options) == 0 {
		return nil
	}
	slices.SortFunc(options, func(a, b *core.Symbol) int { return strings.Compare(a.Name, b.Name) })

	shown := make([]core.Symbol, len(options))
	shownCounts := make([]int, len(options))
	for i, o := range options {
		shown[i] = *o
		shownCounts[i] = counts[o]
	}
	choice := c.ui.AskForActivity(derefCards(hand), shown, shownCounts)
	if choice < 0 {
		return nil
	}
	return options[choice]
}

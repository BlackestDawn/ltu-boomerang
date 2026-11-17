package agent

import (
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/ui"
)

type HumanController struct {
	ui *ui.TextUI
}

func NewHumanController(ui *ui.TextUI) *HumanController {
	return &HumanController{
		ui: ui,
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

func (c *HumanController) ChooseActivity(hand []*core.Card) (*core.Symbol, bool) {
	return &core.Symbol{}, false
}

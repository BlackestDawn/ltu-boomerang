package agent

import (
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/ui"
)

type HumanController struct {
	ui *ui.TextUI
}

func (c *HumanController) ChooseThrow(hand []*core.Card) core.Card {
	return core.Card{}
}

func (c *HumanController) ChooseDraft(hand []*core.Card) core.Card {
	return core.Card{}
}

func (c *HumanController) ChooseBlueSymbol(options []string) (string, bool) {
	return "", false
}

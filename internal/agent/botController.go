package agent

import (
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

type BotController struct {
}

func (c *BotController) ChooseThrow(hand []*core.Card) core.Card {
	return core.Card{}
}

func (c *BotController) ChooseDraft(hand []*core.Card) core.Card {
	return core.Card{}
}

func (c *BotController) ChooseBlueSymbol(options []string) (string, bool) {
	return "", false
}

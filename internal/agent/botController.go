package agent

import (
	"math/rand"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

type BotController struct {
	rng *rand.Rand
}

func NewBotController(rng *rand.Rand) *BotController {
	return &BotController{
		rng: rng,
	}
}

func (c *BotController) ChooseThrow(hand []*core.Card) *core.Card {
	return &core.Card{}
}

func (c *BotController) ChooseDraft(hand []*core.Card) *core.Card {
	return &core.Card{}
}

func (c *BotController) ChooseActivity(hand []*core.Card) (*core.Symbol, bool) {
	return &core.Symbol{}, false
}

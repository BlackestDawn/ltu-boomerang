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
	num := c.rng.Intn(len(hand))
	return hand[num]
}

func (c *BotController) ChooseDraft(hand []*core.Card) *core.Card {
	num := c.rng.Intn(len(hand))
	return hand[num]
}

func (c *BotController) ChooseActivity(hand []*core.Card) (*core.Symbol, bool) {
	activities := core.FilterSymbolCategoryWithCount(hand, core.Activity)
	var retSymbol *core.Symbol
	var top int = 0
	for symbol, num := range activities {
		if num > top {
			top = num
			retSymbol = symbol
		}
	}
	if retSymbol == nil {
		return nil, false
	}
	return retSymbol, true
}

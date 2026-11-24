package agent

import (
	"math/rand"
	"slices"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

type BotController struct {
	rng *rand.Rand
}

// NewBotController creates a bot with its own random source, seeded from rng.
// Bots are asked for their choices concurrently and *rand.Rand is not safe
// for concurrent use, so several bots may be created from one rng but must
// not draw from it while playing.
func NewBotController(rng *rand.Rand) *BotController {
	return &BotController{
		rng: rand.New(rand.NewSource(rng.Int63())),
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

func (c *BotController) ChooseActivity(hand []*core.Card, alreadyScored []*core.Symbol) *core.Symbol {
	activities := core.FilterSymbolCategoryWithCount(hand, core.Activity)
	var retSymbol *core.Symbol
	var top int = 0
	for symbol, num := range activities {
		if slices.Contains(alreadyScored, symbol) {
			continue
		}
		if num > top {
			top = num
			retSymbol = symbol
		}
	}
	return retSymbol
}

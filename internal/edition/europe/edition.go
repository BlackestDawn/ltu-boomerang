package europe

import (
	"slices"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition/shared"
)

type EuropeEditionData struct {
	cards []*core.Card
	rules map[core.Category]core.ScoringRule
}

func (e *EuropeEditionData) Cards() []*core.Card {
	return slices.Clone(e.cards)
}

func (e *EuropeEditionData) Rules() map[core.Category]core.ScoringRule {
	return e.rules
}

// Ties are broken by the passport bonus, then by throw & catch score.
func (e *EuropeEditionData) TieBreaker(all []*core.PlayerState, score int) int {
	return shared.TieBreak(all, score, core.TravelBonus, core.ThrowCatch)
}

func New() *EuropeEditionData {
	return &EuropeEditionData{
		cards: buildCards(),
		rules: buildRules(),
	}
}

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
func (e *EuropeEditionData) GetWinner(sheets []*core.ScoreSheet) int {
	return shared.GetWinner(sheets, core.TravelBonus, core.ThrowCatch)
}

func New() *EuropeEditionData {
	return &EuropeEditionData{
		cards: buildCards(),
		rules: buildRules(),
	}
}

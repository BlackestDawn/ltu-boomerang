package usa

import (
	"slices"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition/shared"
)

type USAEditionData struct {
	cards []*core.Card
	rules map[core.Category]core.ScoringRule
}

func (e *USAEditionData) Cards() []*core.Card {
	return slices.Clone(e.cards)
}

func (e *USAEditionData) Rules() map[core.Category]core.ScoringRule {
	return e.rules
}

// Ties are broken by the coast to coast bonus, then by throw & catch score.
func (e *USAEditionData) GetWinner(sheets []*core.ScoreSheet) int {
	return shared.GetWinner(sheets, core.TravelBonus, core.ThrowCatch)
}

func New() *USAEditionData {
	return &USAEditionData{
		cards: buildCards(),
		rules: buildRules(),
	}
}

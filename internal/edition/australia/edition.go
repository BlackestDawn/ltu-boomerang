package australia

import (
	"slices"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition/shared"
)

type AustraliaEditionData struct {
	cards []*core.Card
	rules map[core.Category]core.ScoringRule
}

func (e *AustraliaEditionData) Cards() []*core.Card {
	return slices.Clone(e.cards)
}

func (e *AustraliaEditionData) Rules() map[core.Category]core.ScoringRule {
	return e.rules
}

// Ties breaking by throw & catch score
func (e *AustraliaEditionData) GetWinner(sheets []*core.ScoreSheet) int {
	return shared.GetWinner(sheets, core.ThrowCatch)
}

func New() *AustraliaEditionData {
	return &AustraliaEditionData{
		cards: buildCards(),
		rules: buildRules(),
	}
}

package australia

import (
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

type AustraliaEditionData struct {
	cards []*core.Card
	rules map[core.Category]core.ScoringRule
}

func (e *AustraliaEditionData) Cards() []*core.Card {
	return e.cards
}

func (e *AustraliaEditionData) Rules() map[core.Category]core.ScoringRule {
	return e.rules
}

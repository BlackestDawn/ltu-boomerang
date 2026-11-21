package australia

import (
	"slices"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
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

// In accordance with official rules, tie breaking is only done through throw & catch score
func (e *AustraliaEditionData) TieBreaker(all []*core.PlayerState, score int) int {
	var topPlayer int
	topTBScore := -1
	for id, player := range all {
		if player.ScoreSheet.Total() == score {
			throwScore := player.ScoreSheet.CategoryTotal(core.ThrowCatch)
			if throwScore > topTBScore {
				topTBScore = throwScore
				topPlayer = id
			}
		}
	}
	return topPlayer
}

func New() *AustraliaEditionData {
	return &AustraliaEditionData{
		cards: buildCards(),
		rules: buildRules(),
	}
}

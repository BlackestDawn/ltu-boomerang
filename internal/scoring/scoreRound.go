package scoring

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

func ScoreRound(round int, states []*core.PlayerState, rules map[core.Category]core.ScoringRule) {
	for i := range states {
		for cat, rule := range rules {
			states[i].ScoreSheet.Add(round, cat, rule.Score(i, states))
		}
	}
}

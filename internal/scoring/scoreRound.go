package scoring

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

func ScoreRound(round int, states []*core.PlayerState, rules map[core.Category]core.ScoringRule) core.RoundScoreEvent {
	result := make(map[int]map[core.Category]int, 0)
	for i := range states {
		result[i] = make(map[core.Category]int)
		for cat, rule := range rules {
			score := rule.Score(i, states)
			states[i].ScoreSheet.Add(round, cat, score)
			result[i][cat] = score
		}
	}
	return result
}

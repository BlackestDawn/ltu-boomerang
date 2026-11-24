package australia

import (
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition/shared"
)

func buildRules() map[core.Category]core.ScoringRule {
	return map[core.Category]core.ScoringRule{
		core.ThrowCatch:  shared.ScoringFunc(shared.ScoreThrowDifference),
		core.Site:        shared.ScoringFunc(shared.ScoreSites),
		core.RegionBonus: shared.ScoringFunc(scoreRegion),
		core.Collection:  shared.ScoringFunc(scoreCollection),
		core.Animal:      shared.ScoringFunc(scoreAnimal),
		core.Activity:    shared.ScoringFunc(scoreActivity),
	}
}

// Every region is worth 3 points to the first player completing it.
var regionBonus = []int{3, 3, 3, 3, 3, 3, 3}

func scoreRegion(p int, all []*core.PlayerState) int {
	return shared.ScoreRegions(p, all, regionClusters, regionBonus)
}

// Collections are doubled when the total is 7 or less.
func scoreCollection(p int, all []*core.PlayerState) int {
	score := shared.SymbolSum(all[p].Draft, core.Collection)
	if score <= 7 {
		score = score * 2
	}
	return score
}

func scoreAnimal(p int, all []*core.PlayerState) int {
	return shared.ScorePairs(p, all, core.Animal)
}

var activityPoints = []int{0, 0, 2, 4, 7, 10, 15}

func scoreActivity(p int, all []*core.PlayerState) int {
	return shared.ScoreActivityTable(p, all, activityPoints)
}

package usa

import (
	"slices"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition/shared"
)

func buildRules() map[core.Category]core.ScoringRule {
	return map[core.Category]core.ScoringRule{
		core.ThrowCatch:  shared.ScoringFunc(scoreThrow),
		core.Site:        shared.ScoringFunc(shared.ScoreSites),
		core.RegionBonus: shared.ScoringFunc(scoreRegion),
		core.Collection:  shared.ScoringFunc(scoreAmericana),
		core.Animal:      shared.ScoringFunc(scoreAnimal),
		core.Activity:    shared.ScoringFunc(scoreActivity),
		core.TravelBonus: shared.ScoringFunc(scoreCoastToCoast),
	}
}

// The Throw value is scored if the Catch is equal or higher.
func scoreThrow(p int, all []*core.PlayerState) int {
	draft := all[p].Draft
	if len(draft) == 0 {
		return 0
	}
	throw, catch := draft[0].Value, draft[len(draft)-1].Value
	if catch >= throw {
		return throw
	}
	return 0
}

func scoreRegion(p int, all []*core.PlayerState) int {
	return shared.ScoreRegions(p, all, regionClusters, regionBonus)
}

// Americana only scores if the total beats the previous round's score. After
// a zero any Americana scores again.
func scoreAmericana(p int, all []*core.PlayerState) int {
	total := shared.SymbolSum(all[p].Draft, core.Collection)
	var previous int
	if earlier := all[p].ScoreSheet.CategoryScores[core.Collection]; len(earlier) > 0 {
		previous = earlier[len(earlier)-1]
	}
	if total > previous {
		return total
	}
	return 0
}

func scoreAnimal(p int, all []*core.PlayerState) int {
	return shared.ScorePairs(p, all, core.Animal)
}

var activityPoints = []int{0, 0, 2, 4, 7, 10, 15}

func scoreActivity(p int, all []*core.PlayerState) int {
	return shared.ScoreActivityTable(p, all, activityPoints)
}

// 7 points for the first to connect the coasts, 3 for the second and, with
// 3 or 4 players, 1 for anyone else.
func scoreCoastToCoast(p int, all []*core.PlayerState) int {
	tiers := []int{7, 3}
	if len(all) >= 3 {
		tiers = append(tiers, 1)
	}
	return shared.RankedBonus(p, all, core.TravelBonus, coastsConnected, tiers, len(all) >= 3)
}

// coastsConnected reports whether the visited towns form an unbroken path
// from a west coast town to an east coast town.
func coastsConnected(visited []byte) bool {
	var queue []byte
	for _, site := range westCoast {
		if slices.Contains(visited, site) {
			queue = append(queue, site)
		}
	}
	reached := slices.Clone(queue)
	for len(queue) > 0 {
		site := queue[0]
		queue = queue[1:]
		if slices.Contains(eastCoast, site) {
			return true
		}
		for _, r := range routes {
			next := r[1]
			switch site {
			case r[0]:
			case r[1]:
				next = r[0]
			default:
				continue
			}
			if slices.Contains(visited, next) && !slices.Contains(reached, next) {
				reached = append(reached, next)
				queue = append(queue, next)
			}
		}
	}
	return false
}

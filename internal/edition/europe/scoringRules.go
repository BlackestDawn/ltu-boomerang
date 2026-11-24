package europe

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
		core.Collection:  shared.ScoringFunc(scoreCuisine),
		core.Animal:      shared.ScoringFunc(scoreTransport),
		core.Activity:    shared.ScoringFunc(scoreTreasure),
		core.TravelBonus: shared.ScoringFunc(scorePassport),
	}
}

// The difference between Throw and Catch.
func scoreThrow(p int, all []*core.PlayerState) int {
	draft := all[p].Draft
	if len(draft) == 0 {
		return 0
	}
	diff := draft[0].Value - draft[len(draft)-1].Value
	if diff < 0 {
		return -diff
	}
	return diff
}

func scoreRegion(p int, all []*core.PlayerState) int {
	return shared.ScoreRegions(p, all, regionClusters, regionBonus)
}

// Cuisine scores its total up to 7; above that, half the total rounded up.
func scoreCuisine(p int, all []*core.PlayerState) int {
	total := shared.SymbolSum(all[p].Draft, core.Collection)
	if total <= 7 {
		return total
	}
	return (total + 1) / 2
}

func scoreTransport(p int, all []*core.PlayerState) int {
	return shared.ScorePairs(p, all, core.Animal)
}

// Points per matching treasure drafted by other players, by player count.
var treasurePoints = map[int]int{2: 4, 3: 3, 4: 2}

// A chosen treasure, held at least once, scores for every matching icon in the
// other players' drafts. Each treasure scores only once per game.
func scoreTreasure(p int, all []*core.PlayerState) int {
	choice := all[p].ChosenActivity
	if choice == nil || slices.Contains(all[p].ScoredActivities, choice) {
		return 0
	}
	if !slices.ContainsFunc(all[p].Draft, func(c *core.Card) bool { return c.Symbols[core.Activity] == choice }) {
		return 0
	}
	var count int
	for id, state := range all {
		if id == p {
			continue
		}
		for _, card := range state.Draft {
			if card.Symbols[core.Activity] == choice {
				count++
			}
		}
	}
	return count * treasurePoints[len(all)]
}

// 7 points for the first to visit a country in every region, 3 for the
// second and, with 3 or 4 players, 1 for the third.
func scorePassport(p int, all []*core.PlayerState) int {
	tiers := []int{7, 3}
	if len(all) >= 3 {
		tiers = append(tiers, 1)
	}
	return shared.RankedBonus(p, all, core.TravelBonus, visitedEveryRegion, tiers, false)
}

func visitedEveryRegion(visited []byte) bool {
	for _, set := range regionClusters {
		if !slices.ContainsFunc(set, func(site byte) bool { return slices.Contains(visited, site) }) {
			return false
		}
	}
	return true
}

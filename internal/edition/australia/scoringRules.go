package australia

import (
	"slices"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

type scoringFunc func(p int, all []*core.PlayerState) int

func (s scoringFunc) Score(p int, a []*core.PlayerState) int { return s(p, a) }

func buildRules() map[core.Category]core.ScoringRule {
	return map[core.Category]core.ScoringRule{
		core.ThrowCatch:  scoringFunc(scoreThrow),
		core.Site:        scoringFunc(scoreSite),
		core.RegionBonus: scoringFunc(scoreRegion),
		core.Collection:  scoringFunc(scoreCollection),
		core.Animal:      scoringFunc(scoreAnimal),
		core.Activity:    scoringFunc(scoreActivity),
	}
}

func scoreThrow(p int, all []*core.PlayerState) int {
	diff := all[p].Draft[0].Value - all[p].Draft[core.HandSize-1].Value
	if diff < 0 {
		return -diff
	}
	return diff
}

func scoreSite(p int, all []*core.PlayerState) int {
	var score int
	for _, card := range all[p].Draft {
		if !slices.Contains(all[p].VisitedSites, card.Site) {
			score++
		}
	}
	return score
}

func scoreRegion(p int, all []*core.PlayerState) int {
	var score int
REGIONSET:
	for _, set := range regionClusters {
		for _, state := range all {
			if core.SliceIsSubset(set, state.VisitedSites) {
				continue REGIONSET
			}
		}
		allVisited := slices.Clone(all[p].VisitedSites)
		for _, card := range all[p].Draft {
			if !slices.Contains(allVisited, card.Site) {
				allVisited = append(allVisited, card.Site)
			}
		}
		if core.SliceIsSubset(set, allVisited) {
			score += 3
		}
	}
	return score
}

func scoreCollection(p int, all []*core.PlayerState) int {
	var score int
	for _, card := range all[p].Draft {
		coll, ok := card.Symbols[core.Collection]
		if ok {
			score += coll.Value
		}
	}
	if score <= 7 {
		score = score * 2
	}
	return score
}

func scoreAnimal(p int, all []*core.PlayerState) int {
	var score int
	animals := core.FilterSymbolCategoryWithCount(all[p].Draft, core.Animal)
	for ani, num := range animals {
		score += ani.Value * (num / 2)
	}
	return score
}

var activityPoints = []int{0, 0, 2, 4, 7, 10, 15}

func scoreActivity(p int, all []*core.PlayerState) int {
	choice := all[p].ChosenActivity
	if choice == nil || slices.Contains(all[p].ScoredActivities, choice) {
		return 0
	}
	var count int
	for _, card := range all[p].Draft {
		if card.Symbols[core.Activity] == choice {
			count++
		}
	}
	return activityPoints[count]
}

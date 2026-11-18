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
		core.TieBreaker:  scoringFunc(scoreTiebreaker),
	}
}

func scoreThrow(p int, all []*core.PlayerState) int {
	diff := all[p].Draft.Cards[0].Value - all[p].Draft.Cards[core.HandSize-1].Value
	if diff < 0 {
		return -diff
	}
	return diff
}

func scoreSite(p int, all []*core.PlayerState) int {
	var score int
	for _, card := range all[p].Draft.Cards {
		if !slices.Contains(all[p].VisitedSites, card.Site) {
			score++
		}
	}
	return score
}

func scoreRegion(p int, all []*core.PlayerState) int {
	var score int

	return score
}

func scoreCollection(p int, all []*core.PlayerState) int {
	var score int
	for _, card := range all[p].Draft.Cards {
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
	animals := core.FilterSymbolCategoryWithCount(all[p].Draft.Cards, core.Animal)
	for ani, num := range animals {
		score += ani.Value * num % 2
	}
	return score
}

func scoreActivity(p int, all []*core.PlayerState) int {
	var score int

	return score
}

func scoreTiebreaker(p int, all []*core.PlayerState) int {
	var playerId int
	var currentHigh int
	for id, player := range all {
		if player.ScoreSheet.CategoryTotal(core.ThrowCatch) > currentHigh {
			playerId = id
		}
	}
	return playerId
}

package engine

import (
	"slices"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

func (g *Game) passHands(round int, catch bool) {
	targets := make([][]*core.Card, len(g.playerStates))
	for i := range g.players {
		var to int
		if !catch {
			to = g.passing.Target(round, i, len(g.playerStates))
		} else {
			to = g.passing.CatchTarget(round, i, len(g.playerStates))
		}
		targets[to] = g.playerStates[i].Hand
	}
	for i := 0; i < len(g.playerStates); i++ {
		g.playerStates[i].Hand = targets[i]
	}
}

func (g *Game) updateEndOfRound() {
	for _, state := range g.playerStates {
		for _, card := range state.Draft {
			if !slices.Contains(state.VisitedSites, card.Site) {
				state.VisitedSites = append(state.VisitedSites, card.Site)
			}
		}
		if state.ChosenActivity != nil {
			state.ScoredActivities = append(state.ScoredActivities, state.ChosenActivity)
		}
	}
}

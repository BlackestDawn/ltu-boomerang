package engine

import (
	"fmt"
	"math/rand"
	"slices"
	"sync"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition"
	"github.com/BlackestDawn/ltu-boomerang/internal/scoring"
)

type Game struct {
	players      []PlayerController
	playerStates []*core.PlayerState
	deck         *core.Deck
	edition      edition.Edition
	events       Listener
	rng          *rand.Rand
	passing      core.PassingStrategy
	numRounds    int
}

func (g *Game) Play() core.Result {
	fmt.Println("Starting game session")
	var wg sync.WaitGroup

	for round := range g.numRounds {
		// (RE-)Initialize deck
		g.deck = &core.Deck{Cards: g.edition.Cards()}
		g.deck.Shuffle(g.rng)
		fmt.Println("Selecting throw card")

		// Deal new hand and reset round-specific data
		for i := range g.players {
			g.playerStates[i].Hand = g.deck.Deal(core.HandSize)
			g.playerStates[i].Draft = g.playerStates[i].Draft[:0]
			g.playerStates[i].ChosenActivity = nil
		}

		// Choose the throw card
		for i := range g.players {
			wg.Go(func() {
				card := g.players[i].ChooseThrow(g.playerStates[i].Hand)
				g.playerStates[i].MoveCardToDraft(card)
			})
		}
		wg.Wait()

		// Draw cards until one remaining
		for len(g.playerStates[0].Hand) > 1 {
			fmt.Println("selecting draft card, remaining in hand:", len(g.playerStates[0].Hand))
			g.PassHands(round)
			for i := range g.players {
				wg.Go(func() {
					card := g.players[i].ChooseDraft(g.playerStates[i].Hand)
					g.playerStates[i].MoveCardToDraft(card)
					g.events.OnEvent(core.Event{
						Type: core.CardsRevelead,
						Data: slices.Clone(g.playerStates[i].Draft[1:]),
					})
				})
			}
			wg.Wait()
		}

		// One card left, pass it around and automatically add to draft as catch card
		// Then calculate the rounds scores
		g.PassHands(round)
		for i := range g.players {
			wg.Go(func() {
				g.playerStates[i].MoveCardToDraft(g.playerStates[i].Hand[0])
				g.events.OnEvent(core.Event{
					Type: core.CardsRevelead,
					Data: slices.Clone(g.playerStates[i].Draft),
				})
				g.playerStates[i].ChosenActivity = g.players[i].ChooseActivity(g.playerStates[i].Draft, g.playerStates[i].ScoredActivities)
			})
		}
		wg.Wait()
		scoring.ScoreRound(round, g.playerStates, g.edition.Rules())

		g.UpdateEndOfRound()
	}

	return scoring.ScoreFinal(g.playerStates, g.edition.TieBreaker)
}

func (g *Game) PassHands(round int) {
	targets := make([][]*core.Card, len(g.playerStates))
	for i := range g.players {
		to := g.passing.Target(round, i, len(g.playerStates))
		targets[to] = g.playerStates[i].Hand
	}
	for i := 0; i < len(g.playerStates); i++ {
		g.playerStates[i].Hand = targets[i]
	}
}

func (g *Game) UpdateEndOfRound() {
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

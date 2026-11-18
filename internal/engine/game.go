package engine

import (
	"fmt"
	"math/rand"
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
}

func (g *Game) Play() Result {
	fmt.Println("Starting game session")
	var wg sync.WaitGroup

	for round := range core.NumRounds { // Initialize deck, hands and choose throw card
		g.deck = &core.Deck{Cards: g.edition.Cards()}
		g.deck.Shuffle(g.rng)
		fmt.Println("Selecting throw card")
		for i := 0; i < len(g.players); i++ {
			g.playerStates[i].Hand = g.deck.Deal(core.HandSize)
		}
		for i := 0; i < len(g.players); i++ {
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
			for i := 0; i < len(g.players); i++ {
				wg.Go(func() {
					card := g.players[i].ChooseDraft(g.playerStates[i].Hand)
					g.playerStates[i].MoveCardToDraft(card)
					g.events.OnEvent(core.Event{
						Type: core.CardsRevelead,
						Data: g.playerStates[i].Draft,
					})
				})
			}
			wg.Wait()
		}

		// One card left, pass it and automatically add to draft
		// then calculate the rounds scores
		g.PassHands(round)
		for i := 0; i < len(g.playerStates); i++ {
			g.playerStates[i].MoveCardToDraft(g.playerStates[i].Hand[0])
			activity, _ := g.players[i].ChooseActivity(g.playerStates[i].Draft)
			scoring.ScoreRound(g.playerStates[i], g.edition.Rules(), activity)
		}
	}

	return Result{
		Ranking: g.playerStates,
		Winner:  0,
	}
}

func (g *Game) PassHands(round int) {
	targets := make([][]*core.Card, len(g.playerStates))
	for i := 0; i < len(g.playerStates); i++ {
		to := g.passing.Target(round, i, len(g.playerStates))
		targets[to] = g.playerStates[i].Hand
	}
	for i := 0; i < len(g.playerStates); i++ {
		g.playerStates[i].Hand = targets[i]
	}
}

package engine

import (
	"fmt"
	"math/rand"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition"
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
	round, seat := 1, 1

	// Round
	g.deck = &core.Deck{Cards: g.edition.Cards()}
	g.deck.Shuffle(g.rng)
	for i := 0; i < len(g.players); i++ {
		fmt.Println("Selecting throw card")
		// TODO: make into go routine
		g.playerStates[i].Hand.Set(g.deck.Deal(core.HandSize))
		card := g.players[i].ChooseThrow(g.playerStates[i].Hand.Cards)
		g.playerStates[i].Hand.Remove(card)
		g.playerStates[i].Draft.Add(card)
	}

	for len(g.playerStates[0].Hand.Cards) > 1 {
		fmt.Println("selecting draft card, remaining in hand:", len(g.playerStates[0].Hand.Cards))
		g.PassHands(round, seat)
		for i := 0; i < len(g.players); i++ {
			// TODO: make into go routine
			card := g.players[i].ChooseDraft(g.playerStates[i].Hand.Cards)
			g.playerStates[i].Hand.Remove(card)
			g.playerStates[i].Draft.Add(card)
		}
		// g.Listener.OnEvent()
	}

	g.PassHands(round, seat)

	return Result{}
}

func (g *Game) PassHands(round, seat int) {}

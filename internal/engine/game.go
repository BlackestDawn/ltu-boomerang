package engine

import (
	"math/rand"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition"
)

type Game struct {
	players      []PlayerController
	playerStates []*core.PlayerState
	deck         *core.Deck
	edition      edition.Edition
	events       core.Publisher
	rng          *rand.Rand
	passing      core.PassingStrategy
	numRounds    int
}

func NewGame(ed edition.Edition, players []PlayerController, rng *rand.Rand, passingDir string, rounds int, events core.Publisher) *Game {
	var states []*core.PlayerState
	for range players {
		states = append(states, core.NewPlayerState(rounds))
	}
	return &Game{
		players:      players,
		playerStates: states,
		edition:      ed,
		events:       events,
		rng:          rng,
		passing:      core.NewPassingStrategy(passingDir),
		numRounds:    rounds,
	}
}

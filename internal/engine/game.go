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
	eventBus     chan core.Event
	rng          *rand.Rand
	passing      core.PassingStrategy
	numRounds    int
}

func NewGame(ed edition.Edition, players []PlayerController, rng *rand.Rand, passingDir string, rounds int, ev chan core.Event) *Game {
	var states []*core.PlayerState
	for range players {
		states = append(states, core.NewPlayerState())
	}
	return &Game{
		players:      players,
		playerStates: states,
		edition:      ed,
		eventBus:     ev,
		rng:          rng,
		passing:      core.NewPassingStrategy(passingDir),
		numRounds:    rounds,
	}
}

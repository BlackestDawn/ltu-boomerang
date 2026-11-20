package engine

import (
	"math/rand"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition"
)

func NewGame(ed edition.Edition, players []PlayerController, rng *rand.Rand, passingDir string, rounds int) *Game {
	var states []*core.PlayerState
	for i := 0; i < len(players); i++ {
		states = append(states, core.NewPlayerState())
	}
	return &Game{
		players:      players,
		playerStates: states,
		edition:      ed,
		rng:          rng,
		passing:      core.NewPassingStrategy(passingDir),
		numRounds:    rounds,
	}
}

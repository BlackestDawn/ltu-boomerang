package scoring

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

type SessionResult struct {
	Ranking []*core.PlayerState
	Winner  int
}

type RoundScore map[int]map[core.Category]int

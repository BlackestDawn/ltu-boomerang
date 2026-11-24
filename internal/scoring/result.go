package scoring

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

type SessionResult struct {
	Scores []*core.ScoreSheet
	Winner int
}

type RoundScore map[int]map[core.Category]int

package scoring

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

type Result struct {
	Ranking []*core.PlayerState
	Winner  int
}

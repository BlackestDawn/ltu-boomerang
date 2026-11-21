package scoring

import (
	"slices"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

func ScoreFinal(all []*core.PlayerState, tieBreaker func(all []*core.PlayerState, score int) int) Result {
	var totals []int
	for _, state := range all {
		totals = append(totals, state.ScoreSheet.Total())
	}
	return Result{
		Ranking: all,
		Winner:  tieBreaker(all, slices.Max(totals)),
	}
}

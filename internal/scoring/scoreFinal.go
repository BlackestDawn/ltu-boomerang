package scoring

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

// ScoreFinal collects the players' score sheets in seat order and lets the
// edition's GetWinner pick the winner among them.
func ScoreFinal(all []*core.PlayerState, getWinner func(sheets []*core.ScoreSheet) int) SessionResult {
	scores := make([]*core.ScoreSheet, len(all))
	for i, state := range all {
		scores[i] = state.ScoreSheet
	}
	return SessionResult{
		Scores: scores,
		Winner: getWinner(scores),
	}
}

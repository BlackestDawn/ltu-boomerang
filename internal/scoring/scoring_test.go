package scoring

import (
	"testing"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

type fixedRule func(p int) int

func (f fixedRule) Score(p int, _ []*core.PlayerState) int { return f(p) }

func newStates(n int) []*core.PlayerState {
	states := make([]*core.PlayerState, n)
	for i := range states {
		states[i] = &core.PlayerState{ScoreSheet: &core.ScoreSheet{CategoryScores: map[core.Category][]int{}}}
	}
	return states
}

// Requirement 10/11: each round every rule is applied to every player, the
// result is returned for display and written to the score sheets.
func TestScoreRoundAppliesAllRules(t *testing.T) {
	states := newStates(3)
	rules := map[core.Category]core.ScoringRule{
		core.Site:   fixedRule(func(p int) int { return p + 1 }),
		core.Animal: fixedRule(func(p int) int { return 10 * p }),
	}

	for round := range 2 {
		result := ScoreRound(round, states, rules)
		if len(result) != 3 {
			t.Fatalf("result for %d players, want 3", len(result))
		}
		for p := range 3 {
			if result[p][core.Site] != p+1 || result[p][core.Animal] != 10*p {
				t.Errorf("round %d player %d: got %v", round, p, result[p])
			}
		}
	}
	for p, s := range states {
		if got := len(s.ScoreSheet.CategoryScores[core.Site]); got != 2 {
			t.Errorf("player %d has %d site entries, want 2", p, got)
		}
		if want := 2*(p+1) + 2*10*p; s.ScoreSheet.Total() != want {
			t.Errorf("player %d total %d, want %d", p, s.ScoreSheet.Total(), want)
		}
	}
}

// Requirement 12: the score sheets are handed to the winner function in seat
// order and its choice is the winner.
func TestScoreFinalUsesGetWinner(t *testing.T) {
	states := newStates(3)

	var got []*core.ScoreSheet
	getWinner := func(sheets []*core.ScoreSheet) int {
		got = sheets
		return 2
	}
	res := ScoreFinal(states, getWinner)
	if res.Winner != 2 {
		t.Errorf("winner %d, want 2", res.Winner)
	}
	if len(got) != 3 || len(res.Scores) != 3 {
		t.Fatalf("getWinner got %d sheets, result has %d, want 3", len(got), len(res.Scores))
	}
	for i, s := range states {
		if got[i] != s.ScoreSheet || res.Scores[i] != s.ScoreSheet {
			t.Errorf("sheet %d is not player %d's", i, i)
		}
	}
}

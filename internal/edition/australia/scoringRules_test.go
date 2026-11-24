package australia

import (
	"testing"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition/shared"
)

// draft builds a player state whose draft is the cards with the given site
// letters, in order.
func draft(t *testing.T, sites string) *core.PlayerState {
	t.Helper()
	bySite := map[byte]*core.Card{}
	for _, c := range buildCards() {
		bySite[c.Site] = c
	}
	p := &core.PlayerState{ScoreSheet: &core.ScoreSheet{CategoryScores: map[core.Category][]int{}}}
	for i := range len(sites) {
		c, ok := bySite[sites[i]]
		if !ok {
			t.Fatalf("no card with site %c", sites[i])
		}
		p.Draft = append(p.Draft, c)
	}
	return p
}

// Requirement 10a: |throw value - catch value|, throw being the first and
// catch the last of the seven drafted cards.
func TestScoreThrowCatch(t *testing.T) {
	cases := []struct {
		sites string
		want  int
	}{
		{"ABCEFGY", 6}, // A=1, Y=7
		{"YBCEFGA", 6}, // order reversed, still positive
		{"EBCDFGH", 0}, // E=4, H=4
		{"UABCDEI", 4}, // U=2, I=6
	}
	for _, c := range cases {
		p := draft(t, c.sites)
		if got := shared.ScoreThrowDifference(0, []*core.PlayerState{p}); got != c.want {
			t.Errorf("%s: got %d, want %d", c.sites, got, c.want)
		}
	}
}

// Requirement 10b: one point per site, but sites visited in earlier rounds
// do not score again.
func TestScoreSite(t *testing.T) {
	p := draft(t, "ABCDEFG")
	if got := shared.ScoreSites(0, []*core.PlayerState{p}); got != 7 {
		t.Errorf("fresh sites: got %d, want 7", got)
	}
	p.VisitedSites = []byte("ABZ")
	if got := shared.ScoreSites(0, []*core.PlayerState{p}); got != 5 {
		t.Errorf("with A, B visited: got %d, want 5", got)
	}
}

// Requirement 10b-i: first player to complete a region gets 3 points.
func TestScoreRegionFirstToComplete(t *testing.T) {
	p0 := draft(t, "ABCDEFG") // completes Western Australia
	p1 := draft(t, "IJKMNOQ") // completes nothing
	all := []*core.PlayerState{p0, p1}

	if got := scoreRegion(0, all); got != 3 {
		t.Errorf("completing player: got %d, want 3", got)
	}
	if got := scoreRegion(1, all); got != 0 {
		t.Errorf("other player: got %d, want 0", got)
	}
}

// Requirement 10b: sites from earlier rounds count towards completing a region.
func TestScoreRegionUsesEarlierVisits(t *testing.T) {
	p := draft(t, "CDEFGIJ")
	p.VisitedSites = []byte("AB")
	if got := scoreRegion(0, []*core.PlayerState{p}); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

// Requirement 10b-ii: players completing the same region in the same round
// all get the bonus.
func TestScoreRegionSimultaneous(t *testing.T) {
	p0 := draft(t, "ABCDEFG")
	p1 := draft(t, "IJKL") // will complete Western Australia via visits
	p1.VisitedSites = []byte("ABCQ")
	p1.Draft = append(p1.Draft, draft(t, "D").Draft...)
	all := []*core.PlayerState{p0, p1}

	for i := range all {
		// p1 completes both Western Australia and Queensland.
		want := 3
		if i == 1 {
			want = 6
		}
		if got := scoreRegion(i, all); got != want {
			t.Errorf("player %d: got %d, want %d", i, got, want)
		}
	}
}

// Requirement 10b-i: a region completed by someone in an earlier round gives
// no bonus to anyone afterwards.
func TestScoreRegionAlreadyClaimed(t *testing.T) {
	p0 := draft(t, "ABCDEFG")
	p1 := draft(t, "IJKMNOQ")
	p1.VisitedSites = []byte("ABCD")
	all := []*core.PlayerState{p0, p1}

	if got := scoreRegion(0, all); got != 0 {
		t.Errorf("late completer: got %d, want 0", got)
	}
	if got := scoreRegion(1, all); got != 0 {
		t.Errorf("earlier completer scored again: got %d, want 0", got)
	}
}

// Requirement 10c: collections summed, doubled when 1-7, as-is when over 7.
func TestScoreCollection(t *testing.T) {
	cases := []struct {
		sites string
		want  int
	}{
		{"BEFHJMN", 0}, // no collections
		{"A", 2},       // Leaves 1 -> doubled
		{"ACD", 12},    // 1+3+2 = 6 -> doubled
		{"KD", 14},     // 5+2 = 7 -> doubled (boundary)
		{"KC", 8},      // 5+3 = 8 -> not doubled (boundary)
		{"KCGWI", 16},  // 5+3+3+3+2 = 16
	}
	for _, c := range cases {
		p := draft(t, c.sites)
		if got := scoreCollection(0, []*core.PlayerState{p}); got != c.want {
			t.Errorf("%s: got %d, want %d", c.sites, got, c.want)
		}
	}
}

// Requirement 10d: points for each matching pair of animals.
func TestScoreAnimal(t *testing.T) {
	cases := []struct {
		sites string
		want  int
	}{
		{"E", 0},     // single emu
		{"BC", 3},    // kangaroo pair
		{"BCJ", 3},   // three kangaroos = one pair
		{"BCJO", 6},  // two kangaroo pairs
		{"ER", 4},    // emu pair
		{"FQ", 5},    // wombat pair
		{"HM", 7},    // koala pair
		{"GX", 9},    // platypus pair
		{"GXHM", 16}, // platypus + koala pairs
		{"BEFH", 0},  // all different
	}
	for _, c := range cases {
		p := draft(t, c.sites)
		if got := scoreAnimal(0, []*core.PlayerState{p}); got != c.want {
			t.Errorf("%s: got %d, want %d", c.sites, got, c.want)
		}
	}
}

// Requirement 10e: activity table 1->0, 2->2, 3->4, 4->7, 5->10, 6->15.
func TestScoreActivityTable(t *testing.T) {
	swim := cardSymbols[swimming]
	swimSites := "HLNSW-" // the six swimming cards
	want := []int{0, 2, 4, 7, 10, 15}
	for n := 1; n <= 6; n++ {
		p := draft(t, swimSites[:n]+"ABCDEFG"[:7-n])
		p.ChosenActivity = swim
		if got := scoreActivity(0, []*core.PlayerState{p}); got != want[n-1] {
			t.Errorf("%d swimming: got %d, want %d", n, got, want[n-1])
		}
	}
}

// Requirement 10e: scoring activities is optional.
func TestScoreActivityNoneChosen(t *testing.T) {
	p := draft(t, "HLNSWAB")
	if got := scoreActivity(0, []*core.PlayerState{p}); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}

// Requirement 10e-ii: an activity scores only once per game.
func TestScoreActivityOncePerGame(t *testing.T) {
	swim := cardSymbols[swimming]
	p := draft(t, "HLNSWAB")
	p.ChosenActivity = swim
	p.ScoredActivities = []*core.Symbol{swim}
	if got := scoreActivity(0, []*core.PlayerState{p}); got != 0 {
		t.Errorf("already scored activity: got %d, want 0", got)
	}

	p.ScoredActivities = []*core.Symbol{cardSymbols[sightseeing]}
	if got := scoreActivity(0, []*core.PlayerState{p}); got != 10 {
		t.Errorf("other activity scored earlier: got %d, want 10", got)
	}
}

// Requirement 10: every scoring category is part of the edition's rules.
func TestRulesCoverAllCategories(t *testing.T) {
	rules := New().Rules()
	for _, cat := range []core.Category{
		core.ThrowCatch, core.Site, core.RegionBonus,
		core.Collection, core.Animal, core.Activity,
	} {
		if _, ok := rules[cat]; !ok {
			t.Errorf("no rule for %s", cat)
		}
	}
}

func sheet(throwCatch, other int) *core.PlayerState {
	return &core.PlayerState{ScoreSheet: &core.ScoreSheet{CategoryScores: map[core.Category][]int{
		core.ThrowCatch: {throwCatch},
		core.Site:       {other},
	}}}
}

// Requirement 12: highest total wins; 12a: ties broken by Throw & Catch.
func TestTieBreaker(t *testing.T) {
	ed := New()
	cases := []struct {
		name string
		all  []*core.PlayerState
		top  int
		want int
	}{
		{"clear winner", []*core.PlayerState{sheet(1, 10), sheet(2, 20), sheet(0, 5)}, 22, 1},
		{"tie, more throw&catch wins", []*core.PlayerState{sheet(2, 18), sheet(9, 11), sheet(1, 3)}, 20, 1},
		{"tie, first beats lower", []*core.PlayerState{sheet(8, 12), sheet(3, 17)}, 20, 0},
		{"non-tied higher T&C ignored", []*core.PlayerState{sheet(15, 0), sheet(1, 19), sheet(2, 18)}, 20, 2},
	}
	for _, c := range cases {
		if got := ed.TieBreaker(c.all, c.top); got != c.want {
			t.Errorf("%s: winner %d, want %d", c.name, got, c.want)
		}
	}
}

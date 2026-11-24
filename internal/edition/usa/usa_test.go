package usa

import (
	"slices"
	"testing"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

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

func TestDeck(t *testing.T) {
	cards := New().Cards()
	if len(cards) != 28 {
		t.Fatalf("deck has %d cards, want 28", len(cards))
	}
	regionOf := map[byte]string{}
	for _, c := range cards {
		if _, dup := regionOf[c.Site]; dup {
			t.Errorf("site %c appears twice", c.Site)
		}
		regionOf[c.Site] = c.Region
		if len(c.Symbols) != 2 {
			t.Errorf("%s has %d symbols, want 2", c.Name, len(c.Symbols))
		}
	}
	for i, cluster := range regionClusters {
		for _, site := range cluster {
			if regionOf[site] != regionNames[i] {
				t.Errorf("site %c is in cluster %s but card says %s", site, regionNames[i], regionOf[site])
			}
		}
	}
	for _, r := range routes {
		if _, ok := regionOf[r[0]]; !ok {
			t.Errorf("route %c-%c: unknown site", r[0], r[1])
		}
		if _, ok := regionOf[r[1]]; !ok {
			t.Errorf("route %c-%c: unknown site", r[0], r[1])
		}
	}
}

func TestScoreThrow(t *testing.T) {
	cases := []struct {
		sites string
		want  int
	}{
		{"MABCDEY", 1}, // 1 -> 7
		{"YABCDEM", 0}, // 7 -> 1
		{"QABCDES", 4}, // 4 -> 4
	}
	for _, c := range cases {
		if got := scoreThrow(0, []*core.PlayerState{draft(t, c.sites)}); got != c.want {
			t.Errorf("%s: got %d, want %d", c.sites, got, c.want)
		}
	}
}

func TestScoreAmericana(t *testing.T) {
	p := draft(t, "DW") // 1 + 5
	if got := scoreAmericana(0, []*core.PlayerState{p}); got != 6 {
		t.Errorf("first round: got %d, want 6", got)
	}
	p.ScoreSheet.CategoryScores[core.Collection] = []int{6}
	if got := scoreAmericana(0, []*core.PlayerState{p}); got != 0 {
		t.Errorf("not more than last round: got %d, want 0", got)
	}
	p.ScoreSheet.CategoryScores[core.Collection] = []int{5}
	if got := scoreAmericana(0, []*core.PlayerState{p}); got != 6 {
		t.Errorf("more than last round: got %d, want 6", got)
	}
	p = draft(t, "D")
	p.ScoreSheet.CategoryScores[core.Collection] = []int{9, 0}
	if got := scoreAmericana(0, []*core.PlayerState{p}); got != 1 {
		t.Errorf("after a zero: got %d, want 1", got)
	}
}

func TestScoreAnimal(t *testing.T) {
	if got := scoreAnimal(0, []*core.PlayerState{draft(t, "ACHNGB")}); got != 3+5+7 {
		t.Errorf("got %d, want 15", got)
	}
}

func TestCoastsConnected(t *testing.T) {
	cases := []struct {
		visited string
		want    bool
	}{
		{"BGHJLKRST", true},  // Portland - Denver - Dallas - Atlanta - Miami
		{"AEMNOUVZ", true},   // Seattle - Great Falls - Chicago - Pittsburgh - NYC
		{"BGHJLKRS", false},  // stops at Atlanta
		{"BGHJLKRSTZ", true}, // extra sites don't matter
		{"AEMNOVZ", false},   // gap at Columbus
		{"TXZ@#", false},     // east coast only
	}
	for _, c := range cases {
		if got := coastsConnected([]byte(c.visited)); got != c.want {
			t.Errorf("%s: got %v, want %v", c.visited, got, c.want)
		}
	}
}

func TestScoreCoastToCoast(t *testing.T) {
	route := "BGHJLKRST"
	newPlayer := func(visited, sites string, earlier ...int) *core.PlayerState {
		p := draft(t, sites)
		p.VisitedSites = []byte(visited)
		p.ScoreSheet.CategoryScores[core.TravelBonus] = slices.Clone(earlier)
		return p
	}

	// Two players finishing in the same round both get 7.
	all := []*core.PlayerState{newPlayer(route[:5], route[5:]), newPlayer(route[:3], route[3:]), newPlayer("", "A")}
	for i, want := range []int{7, 7, 0} {
		if got := scoreCoastToCoast(i, all); got != want {
			t.Errorf("same round player %d: got %d, want %d", i, got, want)
		}
	}

	// Later rounds get 3, then 1 for everyone after that.
	all = []*core.PlayerState{
		newPlayer(route, "A", 0, 7, 0),
		newPlayer(route[:6], route[6:], 0, 0, 0),
		newPlayer(route[:6], route[6:], 0, 0, 0),
	}
	if got := scoreCoastToCoast(0, all); got != 0 {
		t.Errorf("already scored: got %d, want 0", got)
	}
	if got := scoreCoastToCoast(1, all); got != 3 {
		t.Errorf("second: got %d, want 3", got)
	}
	all[0].ScoreSheet.CategoryScores[core.TravelBonus] = []int{7, 0, 0}
	all[2].ScoreSheet.CategoryScores[core.TravelBonus] = []int{0, 3, 0}
	if got := scoreCoastToCoast(1, all); got != 1 {
		t.Errorf("third: got %d, want 1", got)
	}

	// With two players the second to connect gets 3.
	two := []*core.PlayerState{newPlayer(route[:6], route[6:], 0, 0), newPlayer(route, "A", 7, 0)}
	if got := scoreCoastToCoast(0, two); got != 3 {
		t.Errorf("2p second: got %d, want 3", got)
	}
}

func TestGetWinner(t *testing.T) {
	sheet := func(c2c, throw, other int) *core.ScoreSheet {
		return &core.ScoreSheet{CategoryScores: map[core.Category][]int{
			core.TravelBonus: {c2c}, core.ThrowCatch: {throw}, core.Site: {other},
		}}
	}
	ed := New()
	if got := ed.GetWinner([]*core.ScoreSheet{sheet(0, 10, 10), sheet(3, 0, 17)}); got != 1 {
		t.Errorf("coast to coast first: winner %d, want 1", got)
	}
	if got := ed.GetWinner([]*core.ScoreSheet{sheet(3, 1, 16), sheet(3, 5, 12)}); got != 1 {
		t.Errorf("then throw & catch: winner %d, want 1", got)
	}
}

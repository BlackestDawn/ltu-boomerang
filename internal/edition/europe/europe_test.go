package europe

import (
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
	if len(regionBonus) != len(regionClusters) {
		t.Fatalf("%d bonuses for %d regions", len(regionBonus), len(regionClusters))
	}
	count := 0
	for i, cluster := range regionClusters {
		count += len(cluster)
		for _, site := range cluster {
			if regionOf[site] != regionNames[i] {
				t.Errorf("site %c is in cluster %s but card says %s", site, regionNames[i], regionOf[site])
			}
		}
	}
	if count != 28 {
		t.Errorf("regions cover %d sites, want 28", count)
	}
}

func TestScoreCuisine(t *testing.T) {
	cases := []struct {
		sites string
		want  int
	}{
		{"A", 0},
		{"BKNO", 5},  // 1+1+1+2
		{"BTO", 6},   // 1+3+2
		{"BTOM", 4},  // 8 -> 4
		{"TXU", 6},   // 11 -> 6
		{"UTXWV", 8}, // 16 -> 8
	}
	for _, c := range cases {
		if got := scoreCuisine(0, []*core.PlayerState{draft(t, c.sites)}); got != c.want {
			t.Errorf("%s: got %d, want %d", c.sites, got, c.want)
		}
	}
}

func TestScoreTransport(t *testing.T) {
	if got := scoreTransport(0, []*core.PlayerState{draft(t, "HPDGCF@")}); got != 4+3+7 {
		t.Errorf("got %d, want 14", got)
	}
}

func TestScoreTreasure(t *testing.T) {
	artSym := cardSymbols[art]
	me := draft(t, "HBT") // two art icons of its own, which don't count
	me.ChosenActivity = artSym
	// Art: S and T, then E and M, then none.
	others := []*core.PlayerState{draft(t, "TSB"), draft(t, "EMI"), draft(t, "A")}

	for n, want := range map[int]int{2: 2 * 4, 3: 4 * 3, 4: 4 * 2} {
		all := append([]*core.PlayerState{me}, others[:n-1]...)
		if got := scoreTreasure(0, all); got != want {
			t.Errorf("%dp: got %d, want %d", n, got, want)
		}
	}

	all := []*core.PlayerState{me, others[0]}
	me.ScoredActivities = []*core.Symbol{artSym}
	if got := scoreTreasure(0, all); got != 0 {
		t.Errorf("already scored: got %d, want 0", got)
	}
	me.ScoredActivities = nil
	me.ChosenActivity = cardSymbols[music]
	if got := scoreTreasure(0, all); got != 0 {
		t.Errorf("not held: got %d, want 0", got)
	}
}

func TestScoreRegion(t *testing.T) {
	p0 := draft(t, "IJKLM")
	p1 := draft(t, "IJABC")
	p1.VisitedSites = []byte("DE")
	all := []*core.PlayerState{p0, p1}
	if got := scoreRegion(0, all); got != 2+3 {
		t.Errorf("player 0: got %d, want 5", got)
	}
	if got := scoreRegion(1, all); got != 2+4 {
		t.Errorf("player 1: got %d, want 6", got)
	}
}

func TestScorePassport(t *testing.T) {
	every := "AFIKNQU"
	p0 := draft(t, every)
	p1 := draft(t, every[:6])
	all := []*core.PlayerState{p0, p1, draft(t, "B")}
	if got := scorePassport(0, all); got != 7 {
		t.Errorf("first: got %d, want 7", got)
	}
	if got := scorePassport(1, all); got != 0 {
		t.Errorf("incomplete: got %d, want 0", got)
	}

	// Third in a 3 player game gets 1, nobody gets anything after that.
	p0.VisitedSites = []byte(every)
	p0.ScoreSheet.CategoryScores[core.TravelBonus] = []int{7, 0}
	all[2].VisitedSites = []byte(every)
	all[2].ScoreSheet.CategoryScores[core.TravelBonus] = []int{0, 3}
	p1.ScoreSheet.CategoryScores[core.TravelBonus] = []int{0, 0}
	p1.Draft = draft(t, "UB").Draft
	p1.VisitedSites = []byte(every[:6])
	if got := scorePassport(1, all); got != 1 {
		t.Errorf("third: got %d, want 1", got)
	}
	if got := scorePassport(1, all[:2]); got != 3 {
		t.Errorf("second of two: got %d, want 3", got)
	}
}

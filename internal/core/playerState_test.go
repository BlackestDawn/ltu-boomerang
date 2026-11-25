package core

import (
	"slices"
	"testing"
)

func TestMoveCardToDraft(t *testing.T) {
	a, b, c := &Card{Name: "A"}, &Card{Name: "B"}, &Card{Name: "C"}
	p := NewPlayerState(MinRounds)
	p.Hand = []*Card{a, b, c}

	p.MoveCardToDraft(b)
	p.MoveCardToDraft(a)

	if !slices.Equal(p.Hand, []*Card{c}) {
		t.Errorf("hand is %v, want only C", p.Hand)
	}
	// Order matters: the first drafted card is the throw card (requirement 5).
	if !slices.Equal(p.Draft, []*Card{b, a}) {
		t.Errorf("draft is %v, want B then A", p.Draft)
	}
}

func TestScoreSheetTotals(t *testing.T) {
	s := &ScoreSheet{CategoryScores: map[Category][]int{}}
	s.Add(0, Site, 3)
	s.Add(1, Site, 2)
	s.Add(0, Animal, 7)
	s.Add(1, ThrowCatch, 0)

	if got := s.CategoryTotal(Site); got != 5 {
		t.Errorf("site total %d, want 5", got)
	}
	if got := s.CategoryTotal(Collection); got != 0 {
		t.Errorf("empty category total %d, want 0", got)
	}
	if got := s.Total(); got != 12 {
		t.Errorf("total %d, want 12", got)
	}
}

// Add writes by round index. Slices grow only up to the round written, so
// their length is the number of rounds recorded, which scoring rules read.
func TestScoreSheetAddByRound(t *testing.T) {
	s := NewScoreSheet(MinRounds)
	s.Add(0, Site, 3)
	if got := s.CategoryScores[Site]; !slices.Equal(got, []int{3}) {
		t.Errorf("after round 0: %v, want [3]", got)
	}
	s.Add(2, Site, 5)
	if got := s.CategoryScores[Site]; !slices.Equal(got, []int{3, 0, 5}) {
		t.Errorf("skipped round: %v, want [3 0 5]", got)
	}
	s.Add(0, Site, 4)
	if got := s.CategoryScores[Site]; !slices.Equal(got, []int{4, 0, 5}) {
		t.Errorf("rewritten round: %v, want [4 0 5]", got)
	}
	if got := s.CategoryTotal(Site); got != 9 {
		t.Errorf("total %d, want 9", got)
	}

	// A zero-value sheet works too.
	var z ScoreSheet
	z.Add(1, Animal, 2)
	if got := z.CategoryScores[Animal]; !slices.Equal(got, []int{0, 2}) {
		t.Errorf("zero-value sheet: %v, want [0 2]", got)
	}
}

func TestSliceIsSubset(t *testing.T) {
	cases := []struct {
		sub, sup []byte
		want     bool
	}{
		{[]byte("ABCD"), []byte("XDCBAY"), true},
		{[]byte("ABCD"), []byte("ABC"), false},
		{[]byte{}, []byte("A"), true},
		{[]byte("A"), []byte{}, false},
	}
	for _, c := range cases {
		if got := SliceIsSubset(c.sub, c.sup); got != c.want {
			t.Errorf("SliceIsSubset(%q, %q) = %v, want %v", c.sub, c.sup, got, c.want)
		}
	}
}

func TestFilterSymbolCategoryWithCount(t *testing.T) {
	roo, emu := &Symbol{Name: "Kangaroos"}, &Symbol{Name: "Emus"}
	leaf := &Symbol{Name: "Leaves"}
	cards := []*Card{
		{Symbols: map[Category]*Symbol{Animal: roo, Collection: leaf}},
		{Symbols: map[Category]*Symbol{Animal: roo}},
		{Symbols: map[Category]*Symbol{Animal: emu}},
		{Symbols: map[Category]*Symbol{Collection: leaf}},
	}
	got := FilterSymbolCategoryWithCount(cards, Animal)
	if len(got) != 2 || got[roo] != 2 || got[emu] != 1 {
		t.Errorf("got %v, want 2 kangaroos and 1 emu", got)
	}
}

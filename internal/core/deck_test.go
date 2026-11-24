package core

import (
	"math/rand"
	"slices"
	"testing"
)

func numberedDeck(n int) *Deck {
	cards := make([]*Card, n)
	for i := range cards {
		cards[i] = &Card{Name: string(rune('A' + i)), Value: i}
	}
	return &Deck{Cards: cards}
}

// Requirement 3: shuffling reorders the deck without losing or duplicating cards.
func TestDeckShuffleIsPermutation(t *testing.T) {
	d := numberedDeck(28)
	original := slices.Clone(d.Cards)

	d.Shuffle(rand.New(rand.NewSource(42)))

	if len(d.Cards) != len(original) {
		t.Fatalf("deck has %d cards after shuffle, want %d", len(d.Cards), len(original))
	}
	for _, c := range original {
		if !slices.Contains(d.Cards, c) {
			t.Errorf("card %s missing after shuffle", c.Name)
		}
	}
	if slices.Equal(d.Cards, original) {
		t.Error("shuffle left the deck in its original order")
	}
}

// Requirement 3: different seeds give different orders, so the shuffle is
// actually driven by the random source.
func TestDeckShuffleDependsOnRNG(t *testing.T) {
	a, b := numberedDeck(28), numberedDeck(28)
	a.Shuffle(rand.New(rand.NewSource(1)))
	b.Shuffle(rand.New(rand.NewSource(2)))

	same := true
	for i := range a.Cards {
		if a.Cards[i].Name != b.Cards[i].Name {
			same = false
			break
		}
	}
	if same {
		t.Error("two different seeds produced the same order")
	}
}

// Requirement 4: dealing 7 cards to each of 4 players uses the whole 28-card
// deck with no card dealt twice.
func TestDeckDealGivesDisjointHands(t *testing.T) {
	d := numberedDeck(28)
	d.Shuffle(rand.New(rand.NewSource(7)))

	seen := map[*Card]bool{}
	for p := range 4 {
		hand := d.Deal(HandSize)
		if len(hand) != HandSize {
			t.Fatalf("player %d got %d cards, want %d", p, len(hand), HandSize)
		}
		for _, c := range hand {
			if seen[c] {
				t.Errorf("card %s dealt twice", c.Name)
			}
			seen[c] = true
		}
	}
	if len(seen) != 28 {
		t.Errorf("dealt %d distinct cards, want 28", len(seen))
	}
}

// A dealt hand is a copy: changing it must not change the deck.
func TestDeckDealReturnsCopy(t *testing.T) {
	d := numberedDeck(14)
	hand := d.Deal(HandSize)
	hand[0] = nil
	if d.Cards[0] == nil {
		t.Error("modifying a dealt hand modified the deck")
	}
}

// Shuffling resets the deal position so a new round can be dealt.
func TestDeckShuffleResetsDealPosition(t *testing.T) {
	d := numberedDeck(14)
	d.Deal(HandSize)
	d.Deal(HandSize)
	d.Shuffle(rand.New(rand.NewSource(3)))

	hand := d.Deal(HandSize)
	if !slices.Equal(hand, d.Cards[:HandSize]) {
		t.Error("deal after shuffle did not start from the top of the deck")
	}
}

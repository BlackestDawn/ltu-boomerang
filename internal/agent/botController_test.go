package agent

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

func TestBotPicksFromHand(t *testing.T) {
	bot := NewBotController(rand.New(rand.NewSource(1)))
	hand := testHand()
	for range 20 {
		if c := bot.ChooseThrow(hand); !slices.Contains(hand, c) {
			t.Fatalf("throw %v not from hand", c)
		}
		if c := bot.ChooseDraft(hand); !slices.Contains(hand, c) {
			t.Fatalf("draft %v not from hand", c)
		}
	}
}

// Requirement 10e: the bot picks its most common activity that has not been
// scored yet, and nothing when all have been scored.
func TestBotChooseActivity(t *testing.T) {
	bot := NewBotController(rand.New(rand.NewSource(1)))
	hand := testHand() // Swimming x2, Bushwalking x1
	swim := hand[0].Symbols[core.Activity]
	walk := hand[1].Symbols[core.Activity]

	if got := bot.ChooseActivity(hand, nil); got != swim {
		t.Errorf("got %v, want swimming", got)
	}
	if got := bot.ChooseActivity(hand, []*core.Symbol{swim}); got != walk {
		t.Errorf("with swimming scored: got %v, want bushwalking", got)
	}
	if got := bot.ChooseActivity(hand, []*core.Symbol{swim, walk}); got != nil {
		t.Errorf("with all scored: got %v, want nil", got)
	}
}

func TestActivityOptionsExcludeScored(t *testing.T) {
	hand := testHand()
	swim := hand[0].Symbols[core.Activity]

	opts, counts := activityOptions(hand, nil)
	if len(opts) != 2 || opts[0].Name != "Bushwalking" || counts[1] != 2 {
		t.Errorf("got %v %v", opts, counts)
	}
	opts, _ = activityOptions(hand, []*core.Symbol{swim})
	if len(opts) != 1 || opts[0].Name != "Bushwalking" {
		t.Errorf("with swimming scored: got %v", opts)
	}
}

package core

import "testing"

// Requirement 6: by default hands go to the next player, and the last player
// passes to the first.
func TestPassLeftTargetsNextPlayer(t *testing.T) {
	p := NewPassingStrategy("left")
	for n := 2; n <= MaxPlayerCount; n++ {
		for seat := range n {
			want := (seat + 1) % n
			if got := p.Target(0, seat, n); got != want {
				t.Errorf("n=%d seat=%d: passed to %d, want %d", n, seat, got, want)
			}
		}
	}
}

// Unknown strategies fall back to the default (next player).
func TestUnknownStrategyDefaultsToLeft(t *testing.T) {
	p := NewPassingStrategy("")
	if got := p.Target(0, 3, 4); got != 0 {
		t.Errorf("last seat passed to %d, want 0", got)
	}
}

func TestPassRightTargetsPreviousPlayer(t *testing.T) {
	p := NewPassingStrategy("right")
	for n := 2; n <= MaxPlayerCount; n++ {
		for seat := range n {
			want := (seat - 1 + n) % n
			if got := p.Target(0, seat, n); got != want {
				t.Errorf("n=%d seat=%d: passed to %d, want %d", n, seat, got, want)
			}
		}
	}
}

// Requirement 13 (optional rule): alternate between next and previous player
// each round.
func TestPassAlternateSwitchesEachRound(t *testing.T) {
	p := NewPassingStrategy("alternate")
	n := 4
	for round := range 4 {
		for seat := range n {
			want := (seat + 1) % n
			if round%2 == 1 {
				want = (seat - 1 + n) % n
			}
			if got := p.Target(round, seat, n); got != want {
				t.Errorf("round=%d seat=%d: passed to %d, want %d", round, seat, got, want)
			}
		}
	}
}

// Requirement 9: the catch card goes the opposite way of the draft passing,
// i.e. to the previous player when drafting passes to the next player.
func TestCatchTargetIsOppositeOfPassing(t *testing.T) {
	for _, strategy := range []string{"left", "right", "alternate"} {
		t.Run(strategy, func(t *testing.T) {
			p := NewPassingStrategy(strategy)
			for n := 3; n <= MaxPlayerCount; n++ {
				for round := range 4 {
					for seat := range n {
						passTo := p.Target(round, seat, n)
						// Passing to the next player means catching goes to
						// the previous one and vice versa.
						want := (seat - 1 + n) % n
						if passTo == (seat-1+n)%n {
							want = (seat + 1) % n
						}
						if got := p.CatchTarget(round, seat, n); got != want {
							t.Fatalf("n=%d round=%d seat=%d: hands pass to %d, catch card to %d, want %d",
								n, round, seat, passTo, got, want)
						}
					}
				}
			}
		})
	}
}

package engine

import (
	"math/rand"
	"slices"
	"sync"
	"testing"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition/australia"
)

// recorder is a scripted player that always picks a random card and remembers
// every hand it was shown, so tests can follow cards around the table.
type recorder struct {
	rng        *rand.Rand
	throwHands [][]*core.Card   // per round
	throws     []*core.Card     // per round
	draftHands [][][]*core.Card // per round, per pick
	picks      [][]*core.Card   // per round, per pick
	finalDraft [][]*core.Card   // per round, the hand passed to ChooseActivity
	scored     [][]*core.Symbol // per round, alreadyScored passed to ChooseActivity
	chosen     []*core.Symbol   // per round
}

func (r *recorder) ChooseThrow(hand []*core.Card) *core.Card {
	card := hand[r.rng.Intn(len(hand))]
	r.throwHands = append(r.throwHands, slices.Clone(hand))
	r.throws = append(r.throws, card)
	r.draftHands = append(r.draftHands, nil)
	r.picks = append(r.picks, nil)
	return card
}

func (r *recorder) ChooseDraft(hand []*core.Card) *core.Card {
	card := hand[r.rng.Intn(len(hand))]
	last := len(r.draftHands) - 1
	r.draftHands[last] = append(r.draftHands[last], slices.Clone(hand))
	r.picks[last] = append(r.picks[last], card)
	return card
}

func (r *recorder) ChooseActivity(hand []*core.Card, alreadyScored []*core.Symbol) *core.Symbol {
	r.finalDraft = append(r.finalDraft, slices.Clone(hand))
	r.scored = append(r.scored, slices.Clone(alreadyScored))
	var choice *core.Symbol
	for _, c := range hand {
		if s, ok := c.Symbols[core.Activity]; ok && !slices.Contains(alreadyScored, s) {
			choice = s
			break
		}
	}
	r.chosen = append(r.chosen, choice)
	return choice
}

// eventLog collects published events; Publish is called from several goroutines.
type eventLog struct {
	mu     sync.Mutex
	events []core.Event
}

func (l *eventLog) Publish(ev core.Event) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = append(l.events, ev)
}

// rounds splits the log at each RoundBeginsEvent.
func (l *eventLog) rounds() [][]core.Event {
	var out [][]core.Event
	for _, ev := range l.events {
		if _, ok := ev.(core.RoundBeginsEvent); ok {
			out = append(out, nil)
			continue
		}
		out[len(out)-1] = append(out[len(out)-1], ev)
	}
	return out
}

func playGame(t *testing.T, n int, passing string, seed int64) (*Game, []*recorder, *eventLog) {
	t.Helper()
	recs := make([]*recorder, n)
	controllers := make([]PlayerController, n)
	for i := range n {
		recs[i] = &recorder{rng: rand.New(rand.NewSource(seed + int64(i)))}
		controllers[i] = recs[i]
	}
	log := &eventLog{}
	g := NewGame(australia.New(), controllers, rand.New(rand.NewSource(seed)), passing, core.MinRounds, log)
	g.Play()
	return g, recs, log
}

// without returns the cards of hand except card.
func without(hand []*core.Card, card *core.Card) []*core.Card {
	return slices.DeleteFunc(slices.Clone(hand), func(c *core.Card) bool { return c == card })
}

func sameCards(a, b []*core.Card) bool {
	if len(a) != len(b) {
		return false
	}
	for _, c := range a {
		if !slices.Contains(b, c) {
			return false
		}
	}
	return true
}

func names(cards []*core.Card) []string {
	out := make([]string, len(cards))
	for i, c := range cards {
		out[i] = c.Name
	}
	return out
}

// Requirement 1: games with 2, 3 and 4 players run to completion.
func TestPlayWithTwoToFourPlayers(t *testing.T) {
	for n := 2; n <= 4; n++ {
		g, _, _ := playGame(t, n, "left", 1)
		if len(g.playerStates) != n {
			t.Errorf("game with %d players has %d states", n, len(g.playerStates))
		}
	}
}

// Requirement 3 and 4: every round each player is dealt 7 cards from a freshly
// shuffled 28-card deck, with no card in two hands.
func TestDealSevenUniqueCardsEachRound(t *testing.T) {
	for n := 2; n <= 4; n++ {
		_, recs, _ := playGame(t, n, "left", int64(n))
		for round := range core.MinRounds {
			seen := map[*core.Card]bool{}
			for p, r := range recs {
				hand := r.throwHands[round]
				if len(hand) != core.HandSize {
					t.Errorf("n=%d round %d player %d dealt %d cards, want 7", n, round, p, len(hand))
				}
				for _, c := range hand {
					if seen[c] {
						t.Errorf("n=%d round %d: %s dealt twice", n, round, c.Name)
					}
					seen[c] = true
				}
			}
		}
		// The deck is reshuffled between rounds, so player 0's first hand
		// should not be the same in every round.
		if sameCards(recs[0].throwHands[0], recs[0].throwHands[1]) &&
			sameCards(recs[0].throwHands[1], recs[0].throwHands[2]) {
			t.Errorf("n=%d: player 0 got the same hand every round", n)
		}
	}
}

// Requirement 5: the throw card is kept as the first card of the draft and is
// not revealed to the other players before the round ends.
func TestThrowCardIsHiddenUntilRoundEnd(t *testing.T) {
	_, recs, log := playGame(t, 3, "left", 5)
	for round, events := range log.rounds() {
		for _, ev := range events {
			if rc, ok := ev.(core.RevealCardEvent); ok {
				if rc.Card.Name == recs[rc.Player].throws[round].Name {
					t.Errorf("round %d: throw card %s of player %d revealed during draft",
						round, rc.Card.Name, rc.Player)
				}
			}
		}
		for p, r := range recs {
			if r.finalDraft[round][0] != r.throws[round] {
				t.Errorf("round %d player %d: first draft card is %s, want throw card %s",
					round, p, r.finalDraft[round][0].Name, r.throws[round].Name)
			}
		}
	}
}

// Requirement 6: after each pick the remaining hand goes to the next player,
// the last player passing to the first.
func TestHandsPassToNextPlayer(t *testing.T) {
	for n := 2; n <= 4; n++ {
		_, recs, _ := playGame(t, n, "left", 10+int64(n))
		checkPassing(t, recs, func(round, seat int) int { return (seat + 1) % n })
	}
}

// Requirement 13 (optional rule): alternate passing direction each round.
func TestHandsPassAlternating(t *testing.T) {
	for n := 3; n <= 4; n++ {
		_, recs, _ := playGame(t, n, "alternate", 20+int64(n))
		checkPassing(t, recs, func(round, seat int) int {
			if round%2 == 0 {
				return (seat + 1) % n
			}
			return (seat - 1 + n) % n
		})
	}
}

func checkPassing(t *testing.T, recs []*recorder, next func(round, seat int) int) {
	t.Helper()
	n := len(recs)
	for round := range core.MinRounds {
		for seat, r := range recs {
			to := recs[next(round, seat)]
			passed := without(r.throwHands[round], r.throws[round])
			if !sameCards(passed, to.draftHands[round][0]) {
				t.Errorf("n=%d round %d: player %d passed %v after throw, next player got %v",
					n, round, seat, names(passed), names(to.draftHands[round][0]))
			}
			for pick := 0; pick < len(r.draftHands[round])-1; pick++ {
				passed := without(r.draftHands[round][pick], r.picks[round][pick])
				if !sameCards(passed, to.draftHands[round][pick+1]) {
					t.Errorf("n=%d round %d pick %d: player %d passed %v, next player got %v",
						n, round, pick, seat, names(passed), names(to.draftHands[round][pick+1]))
				}
			}
		}
	}
}

// Requirement 7 and 8: each player picks and shows one card per pass until one
// card remains, i.e. 5 shown picks per player per round, then all drafted
// cards are shown at the end of the round.
func TestDraftedCardsAreShown(t *testing.T) {
	n := 4
	_, recs, log := playGame(t, n, "left", 3)
	for round, events := range log.rounds() {
		shown := map[int][]string{}
		drafts := map[int]core.RevealDraftEvent{}
		for _, ev := range events {
			switch e := ev.(type) {
			case core.RevealCardEvent:
				shown[e.Player] = append(shown[e.Player], e.Card.Name)
			case core.RevealDraftEvent:
				drafts[e.Player] = e
			}
		}
		for p, r := range recs {
			if len(r.picks[round]) != core.HandSize-2 {
				t.Errorf("round %d player %d made %d picks, want 5", round, p, len(r.picks[round]))
			}
			slices.Sort(shown[p])
			want := names(r.picks[round])
			slices.Sort(want)
			if !slices.Equal(shown[p], want) {
				t.Errorf("round %d player %d: shown %v, want picks %v", round, p, shown[p], want)
			}
			d, ok := drafts[p]
			if !ok {
				t.Errorf("round %d: draft of player %d never shown", round, p)
				continue
			}
			if len(d.Cards) != core.HandSize {
				t.Errorf("round %d player %d: %d cards shown at end, want 7", round, p, len(d.Cards))
			}
		}
	}
}

// Requirement 9: the last card is passed to the previous player (the first
// player passes it to the last) and becomes that player's catch card.
func TestCatchCardGoesToPreviousPlayer(t *testing.T) {
	for n := 2; n <= 4; n++ {
		_, recs, _ := playGame(t, n, "left", 30+int64(n))
		for round := range core.MinRounds {
			for seat, r := range recs {
				hands := r.draftHands[round]
				last := without(hands[len(hands)-1], r.picks[round][len(hands)-1])
				if len(last) != 1 {
					t.Fatalf("round %d player %d has %d cards left, want 1", round, seat, len(last))
				}
				prev := recs[(seat-1+n)%n]
				catch := prev.finalDraft[round][core.HandSize-1]
				if catch != last[0] {
					t.Errorf("n=%d round %d: player %d's last card %s should be caught by player %d, who caught %s",
						n, round, seat, last[0].Name, (seat-1+n)%n, catch.Name)
				}
			}
		}
	}
}

// Requirement 10b: visited sites are remembered between rounds.
// Requirement 10e-ii: chosen activities are remembered so they can't be
// scored again.
func TestUpdateEndOfRound(t *testing.T) {
	cards := australia.New().Cards()
	swim := &core.Symbol{Name: "Swimming"}
	state := core.NewPlayerState()
	state.VisitedSites = []byte{cards[0].Site}
	state.Draft = cards[:3]
	state.ChosenActivity = swim
	g := &Game{playerStates: []*core.PlayerState{state}}

	g.UpdateEndOfRound()

	if !slices.Equal(state.VisitedSites, []byte{cards[0].Site, cards[1].Site, cards[2].Site}) {
		t.Errorf("visited sites %q, want each site once", state.VisitedSites)
	}
	if !slices.Equal(state.ScoredActivities, []*core.Symbol{swim}) {
		t.Errorf("scored activities %v, want [swimming]", state.ScoredActivities)
	}
}

// Requirement 10e-ii: activities already scored are passed to the player, and
// one choice per round is recorded.
func TestScoredActivitiesCarryOver(t *testing.T) {
	_, recs, _ := playGame(t, 2, "left", 8)
	for p, r := range recs {
		var want []*core.Symbol
		for round := range core.MinRounds {
			if !slices.Equal(r.scored[round], want) {
				t.Errorf("player %d round %d: told %v already scored, want %v", p, round, r.scored[round], want)
			}
			if r.chosen[round] != nil {
				want = append(want, r.chosen[round])
			}
		}
	}
}

// Requirement 11 and 12: four rounds are played, each round is scored and
// the scores shown, and every category gets one entry per round.
func TestFourRoundsScored(t *testing.T) {
	g, _, log := playGame(t, 3, "left", 9)

	var begins, scores int
	for _, ev := range log.events {
		switch e := ev.(type) {
		case core.RoundBeginsEvent:
			if int(e) != begins {
				t.Errorf("round %d announced as %d", begins, e)
			}
			begins++
		case core.RoundScoreEvent:
			scores++
			if len(e) != 3 {
				t.Errorf("round score for %d players, want 3", len(e))
			}
		}
	}
	if begins != 4 || scores != 4 {
		t.Errorf("%d rounds started and %d scored, want 4 each", begins, scores)
	}
	for p, s := range g.playerStates {
		for cat := range g.edition.Rules() {
			if got := len(s.ScoreSheet.CategoryScores[cat]); got != 4 {
				t.Errorf("player %d has %d %s entries, want 4", p, got, cat)
			}
		}
	}
}

// Requirement 12: the reported winner has the highest total, ties broken by
// Throw & Catch.
func TestPlayReportsWinner(t *testing.T) {
	for seed := range int64(10) {
		recs := []PlayerController{
			&recorder{rng: rand.New(rand.NewSource(seed))},
			&recorder{rng: rand.New(rand.NewSource(seed + 100))},
			&recorder{rng: rand.New(rand.NewSource(seed + 200))},
		}
		g := NewGame(australia.New(), recs, rand.New(rand.NewSource(seed)), "left", core.MinRounds, &eventLog{})
		res := g.Play()

		w := res.Ranking[res.Winner].ScoreSheet
		for p, s := range res.Ranking {
			if s.ScoreSheet.Total() > w.Total() {
				t.Errorf("seed %d: player %d has %d, more than winner's %d", seed, p, s.ScoreSheet.Total(), w.Total())
			}
			if s.ScoreSheet.Total() == w.Total() &&
				s.ScoreSheet.CategoryTotal(core.ThrowCatch) > w.CategoryTotal(core.ThrowCatch) {
				t.Errorf("seed %d: tie lost by player %d with more Throw & Catch", seed, p)
			}
		}
	}
}

package engine

import (
	"sync"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/scoring"
)

func (g *Game) Play() scoring.SessionResult {
	var wg sync.WaitGroup

	for round := range g.numRounds {
		// (RE-)Initialize deck
		g.deck = &core.Deck{Cards: g.edition.Cards()}
		g.deck.Shuffle(g.rng)

		// Deal new hand and reset round-specific data
		for i := range g.players {
			g.playerStates[i].Hand = g.deck.Deal(core.HandSize)
			g.playerStates[i].Draft = g.playerStates[i].Draft[:0]
			g.playerStates[i].ChosenActivity = nil
		}
		g.events.Publish(core.RoundBeginsEvent(round))

		// Choose the throw card
		for i := range g.players {
			wg.Go(func() {
				card := g.players[i].ChooseThrow(g.playerStates[i].Hand)
				g.playerStates[i].MoveCardToDraft(card)
			})
		}
		wg.Wait()

		// Draw cards until one remaining
		for len(g.playerStates[0].Hand) > 1 {
			g.PassHands(round, false)
			for i := range g.players {
				wg.Go(func() {
					card := g.players[i].ChooseDraft(g.playerStates[i].Hand)
					g.playerStates[i].MoveCardToDraft(card)
					g.events.Publish(core.RevealCardEvent{Player: i, Card: *card})
				})
			}
			wg.Wait()
		}

		// One card left, pass it around and automatically add to draft as catch card
		// Then calculate the rounds scores
		g.PassHands(round, true)
		for i := range g.players {
			wg.Go(func() {
				g.playerStates[i].MoveCardToDraft(g.playerStates[i].Hand[0])
				g.events.Publish(core.RevealDraftEvent{Player: i, Cards: core.DerefSlice(g.playerStates[i].Draft)})
				g.playerStates[i].ChosenActivity = g.players[i].ChooseActivity(g.playerStates[i].Draft, g.playerStates[i].ScoredActivities)
			})
		}
		wg.Wait()
		roundScore := scoring.ScoreRound(round, g.playerStates, g.edition.Rules())
		g.events.Publish(core.RoundScoreEvent(roundScore))

		g.UpdateEndOfRound()
	}

	return scoring.ScoreFinal(g.playerStates, g.edition.GetWinner)
}

package ui

import (
	"slices"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/scoring"
)

// OnEvent updates the UI's view of the game state. It may be called from a
// different goroutine than the one rendering, so state access is locked.
func (ui *TextUI) OnEvent(e core.Event) {
	switch ev := e.(type) {
	case core.RevealCardEvent:
		ui.mu.Lock()
		ui.draftCards[ev.Player] = append(ui.draftCards[ev.Player], ev.Card)
		ui.mu.Unlock()
	case core.RevealDraftEvent:
		ui.mu.Lock()
		ui.draftCards[ev.Player] = slices.Clone(ev.Cards)
		ui.mu.Unlock()
	case core.RoundScoreEvent:
		ui.ShowRoundResult(scoring.RoundScore(ev))
	case core.RoundBeginsEvent:
		ui.mu.Lock()
		ui.round = int(ev)
		for p := range ui.draftCards {
			ui.draftCards[p] = ui.draftCards[p][:0]
		}
		ui.mu.Unlock()
	}
}

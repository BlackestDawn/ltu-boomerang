package agent

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

func derefCards(cards []*core.Card) []core.Card {
	out := make([]core.Card, len(cards))
	for i, c := range cards {
		out[i] = *c
	}
	return out
}

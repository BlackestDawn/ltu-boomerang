package engine

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

type Listener interface {
	OnEvent(e core.Event)
}

type PlayerController interface {
	ChooseThrow(hand []*core.Card) *core.Card
	ChooseDraft(hand []*core.Card) *core.Card
	ChooseActivity(hand []*core.Card, alreadyScored []*core.Symbol) *core.Symbol
}

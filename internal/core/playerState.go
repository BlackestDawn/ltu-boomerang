package core

import (
	"slices"
)

type PlayerState struct {
	Hand             []*Card
	Draft            []*Card
	VisitedSites     []byte
	ChosenActivity   *Symbol
	ScoredActivities []*Symbol
	ScoreSheet       *ScoreSheet
}

func NewPlayerState(rounds int) *PlayerState {
	return &PlayerState{
		Hand:             []*Card{},
		Draft:            []*Card{},
		VisitedSites:     make([]byte, 0),
		ChosenActivity:   nil,
		ScoredActivities: make([]*Symbol, 0),
		ScoreSheet:       NewScoreSheet(rounds),
	}
}

func (p *PlayerState) MoveCardToDraft(card *Card) {
	p.Draft = append(p.Draft, card)
	p.Hand = slices.DeleteFunc(p.Hand, func(c *Card) bool { return c == card })
}

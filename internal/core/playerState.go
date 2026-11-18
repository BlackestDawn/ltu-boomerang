package core

import (
	"fmt"
	"slices"
)

type PlayerState struct {
	Hand             []*Card
	Draft            []*Card
	VisitedSites     []byte
	ScoredActivities []string
	ScoreSheet       *ScoreSheet
}

func NewPlayerState() *PlayerState {
	fmt.Println("Initilizing new player state")
	return &PlayerState{
		Hand:             []*Card{},
		Draft:            []*Card{},
		VisitedSites:     make([]byte, 0),
		ScoredActivities: make([]string, 0),
		ScoreSheet: &ScoreSheet{
			CategoryScores: make(map[Category][]int),
		},
	}
}

func (p *PlayerState) MoveCardToDraft(card *Card) {
	p.Draft = append(p.Draft, card)
	p.Hand = slices.DeleteFunc(p.Hand, func(c *Card) bool { return c == card })
}

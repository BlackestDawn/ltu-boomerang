package core

import "fmt"

type PlayerState struct {
	Hand             CardPile
	Draft            CardPile
	VisitedSites     []byte
	ScoredActivities []string
	ScoreSheet       *ScoreSheet
}

func NewPlayerState() *PlayerState {
	fmt.Println("Initilizing new player state")
	return &PlayerState{
		Hand: CardPile{
			Cards: []*Card{},
		},
		Draft: CardPile{
			Cards: []*Card{},
		},
		VisitedSites:     make([]byte, 0),
		ScoredActivities: make([]string, 0),
		ScoreSheet: &ScoreSheet{
			CategoryScores: make(map[Category][]int),
		},
	}
}

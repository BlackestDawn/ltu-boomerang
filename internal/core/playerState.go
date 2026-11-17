package core

type PlayerState struct {
	Hand             CardPile
	Draft            CardPile
	VisitedSites     []byte
	ScoredActivities []string
	ScoreSheet       *ScoreSheet
}

package core

type PlayerState struct {
	Hand             []*Card
	Draft            []*Card
	VisitedSites     []byte
	ScoredActivities []string
	ScoreSheet       *ScoreSheet
}

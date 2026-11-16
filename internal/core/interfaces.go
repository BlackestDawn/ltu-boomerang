package core

type ScoringRule interface {
	Score(p int, all []*PlayerState) int
}

package core

type ScoringRule interface {
	Score(p int, all []*PlayerState) int
}

type PassingStrategy interface {
	Target(round, seat, n int) int
	CatchTarget(round, seat, n int) int
}

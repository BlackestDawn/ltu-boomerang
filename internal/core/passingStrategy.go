package core

type Passing struct {
	strategy string
}

func NewPassing(strategy string) Passing {
	return Passing{
		strategy: strategy,
	}
}

func (p Passing) Target(round, seat, n int) int {
	return 0
}

func (p Passing) CatchTarget(round, seat, n int) int {
	return 0
}

package core

type Passing struct {
	target func(round, seat, n int) int
	catch  func(round, seat, n int) int
}

func NewPassingStrategy(strategy string) PassingStrategy {
	switch strategy {
	case "right":
		return Passing{
			target: passRight,
			catch:  passRight,
		}
	case "alternate":
		return Passing{
			target: alternate,
			catch:  alternate,
		}
	default:
		return Passing{
			target: passLeft,
			catch:  passLeft,
		}
	}
}

func (p Passing) Target(round, seat, n int) int {
	return p.target(round, seat, n)
}

func (p Passing) CatchTarget(round, seat, n int) int {
	return p.catch(round, seat, n)
}

func passLeft(round, seat, n int) int {
	target := seat + 1
	if target >= n {
		return 0
	}
	return target
}

func passRight(round, seat, n int) int {
	target := seat - 1
	if target < 0 {
		return n - 1
	}
	return target
}

func alternate(round, seat, n int) int {
	if round%2 == 0 {
		return passLeft(round, seat, n)
	}
	return passRight(round, seat, n)
}

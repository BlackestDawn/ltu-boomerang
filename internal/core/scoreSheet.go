package core

type ScoreSheet struct {
	categoryScores map[Category][]int
}

func (s *ScoreSheet) Add(round int, category Category, points int) {}

func (s *ScoreSheet) CategoryTotal(category Category) int {
	var total int
	for _, val := range s.categoryScores[category] {
		total += val
	}
	return total
}

func (s *ScoreSheet) Total() int {
	return 0
}

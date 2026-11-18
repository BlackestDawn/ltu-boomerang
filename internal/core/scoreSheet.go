package core

type ScoreSheet struct {
	CategoryScores map[Category][]int
}

func (s *ScoreSheet) Add(round int, category Category, points int) {
	s.CategoryScores[category] = append(s.CategoryScores[category], points)
}

func (s *ScoreSheet) CategoryTotal(category Category) int {
	var total int
	for _, val := range s.CategoryScores[category] {
		total += val
	}
	return total
}

func (s *ScoreSheet) Total() int {
	var total int
	for cat := range s.CategoryScores {
		total += s.CategoryTotal(cat)
	}
	return total
}

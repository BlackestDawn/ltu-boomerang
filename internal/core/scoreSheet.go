package core

// ScoreSheet holds a player's points per category, indexed by round. The
// length of a category's slice is the number of rounds recorded for it so
// far; scoring rules rely on that to find the previous rounds' entries.
type ScoreSheet struct {
	maxRounds      int // capacity hint for new category slices
	CategoryScores map[Category][]int
}

// NewScoreSheet creates an empty score sheet for a game of rounds rounds.
func NewScoreSheet(rounds int) *ScoreSheet {
	return &ScoreSheet{
		maxRounds:      rounds,
		CategoryScores: make(map[Category][]int),
	}
}

// Add records points for category in round (0-based), replacing any earlier
// entry for that round. Skipped rounds are filled with 0. The slice only grows
// up to round, never ahead of it, so its length stays the rounds recorded.
func (s *ScoreSheet) Add(round int, category Category, points int) {
	if s.CategoryScores == nil {
		s.CategoryScores = make(map[Category][]int)
	}
	scores := s.CategoryScores[category]
	if scores == nil {
		scores = make([]int, 0, max(s.maxRounds, round+1))
	}
	for len(scores) <= round {
		scores = append(scores, 0)
	}
	scores[round] = points
	s.CategoryScores[category] = scores
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

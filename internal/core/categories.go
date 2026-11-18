package core

// categories used for scoring and symbols
type Category int

const (
	ThrowCatch Category = iota
	Site
	RegionBonus
	Collection // Scoring by total sum, green symbol
	Animal     // Scoring for each pair, yellow symbol
	Activity   // Optional scoring once per icon, blue symbol
	TieBreaker
)

var CatStrings map[Category]string

func (c Category) String() string {
	return CatStrings[c]
}

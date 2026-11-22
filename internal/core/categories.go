package core

import "fmt"

// categories used for scoring and symbols
type Category int

const (
	ThrowCatch Category = iota
	Site
	RegionBonus
	Collection // Scoring by total sum, green symbol
	Animal     // Scoring for each pair, yellow symbol
	Activity   // Optional scoring once per icon, blue symbol
)

var CatStrings map[Category]string

func (c Category) String() string {
	if s, ok := CatStrings[c]; ok {
		return s
	}
	if s, ok := DefaultCatStrings[c]; ok {
		return s
	}
	return fmt.Sprintf("Category %d", int(c))
}

package core

// Base struct for card symbols
type Symbol struct {
	Name  string
	Value int
}

func FilterSymbolCategoryWithCount(cards []*Card, cat Category) map[*Symbol]int {
	symbols := make(map[*Symbol]int)
	for _, card := range cards {
		if symbol, ok := card.Symbols[cat]; ok {
			symbols[symbol]++
		}
	}
	return symbols
}

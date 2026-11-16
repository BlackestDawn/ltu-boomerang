package core

// Symbol categories matching the color one the printed cards
type SymbolCategory int

const (
	Green  SymbolCategory = iota // Scoring as total sum
	Yellow                       // Scoring only if having two or more
	Blue                         // Optional scoring once per round
)

// Base struct for symbols
type Symbol struct {
	Name  string
	value int
}

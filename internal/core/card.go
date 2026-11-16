package core

type Card struct {
	Name    string
	Site    byte
	Region  string
	Value   int
	Symbols map[SymbolCategory]Symbol
}

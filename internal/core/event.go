package core

type Event struct {
	Type EventType `json:"type"`
	Data any       `json:"data"`
}

type EventType int

const (
	CardsRevelead EventType = iota
	RoundScored
)

func (e EventType) String() string {
	return [...]string{
		"cardsRevelead",
		"roundScored",
	}[e]
}

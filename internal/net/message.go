package net

import (
	"encoding/json"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

type MessageType int

const (
	// Host -> client
	Welcome MessageType = iota
	Event
	AskThrow
	AskDraw
	AskActivity
	Final
	Error
	Bye
	// Client -> host
	Hello
	Answer
)

func (t MessageType) String() string {
	switch t {
	case Welcome:
		return "welcome"
	case Event:
		return "event"
	case AskThrow:
		return "askThrow"
	case AskDraw:
		return "askDraw"
	case AskActivity:
		return "askActivity"
	case Final:
		return "final"
	case Error:
		return "error"
	case Bye:
		return "bye"
	case Hello:
		return "hello"
	case Answer:
		return "answer"
	}
	return "unknown"
}

// Message is what goes over the wire, one JSON object per line. Data holds the
// payload for Type still encoded, use Decode to read it into the matching
// payload struct below.
type Message struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// NewMessage builds a message with payload encoded as its data. A nil payload
// gives a message without data.
func NewMessage(t MessageType, payload any) (Message, error) {
	if payload == nil {
		return Message{Type: t}, nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return Message{}, err
	}
	return Message{Type: t, Data: data}, nil
}

// NewEventMessage wraps a game event, using the core event envelope as data.
func NewEventMessage(ev core.Event) (Message, error) {
	data, err := core.MarshalEvent(ev)
	if err != nil {
		return Message{}, err
	}
	return Message{Type: Event, Data: data}, nil
}

// Decode reads the message data into v.
func (m Message) Decode(v any) error {
	return json.Unmarshal(m.Data, v)
}

// Payloads, named after the message type they belong to.

type HelloData struct {
	Version int `json:"version"`
}

type WelcomeData struct {
	Seat       int `json:"seat"`
	NumPlayers int `json:"numPlayers"`
	Rounds     int `json:"rounds"`
}

// AskCardData is the payload of both AskThrow and AskDraw.
type AskCardData struct {
	Hand []core.Card `json:"hand"`
}

type AskActivityData struct {
	Draft   []core.Card   `json:"draft"`
	Options []core.Symbol `json:"options"`
	Counts  []int         `json:"counts"`
}

// AnswerData is the index of the chosen option in the ask it answers, -1
// meaning skip where that is allowed.
type AnswerData struct {
	Choice int `json:"choice"`
}

// FinalData holds every player's score sheet, indexed by seat.
type FinalData struct {
	ScoreSheets []core.ScoreSheet `json:"scoreSheets"`
	Winner      int               `json:"winner"`
}

type ErrorData struct {
	Message string `json:"message"`
}

// ProtocolVersion is sent in Hello so mismatched builds can be refused.
const ProtocolVersion = 1

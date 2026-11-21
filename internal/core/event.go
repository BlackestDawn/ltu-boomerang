package core

import (
	"encoding/json"
	"fmt"
)

type Event interface {
	isEvent()
}

type RevealCardEvent struct {
	Player int  `json:"player"`
	Card   Card `json:"card"`
}

type RoundScoreEvent map[int]map[Category]int

type RoundBeginsEvent int

func (RevealCardEvent) isEvent()  {}
func (RoundScoreEvent) isEvent()  {}
func (RoundBeginsEvent) isEvent() {}

type envelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func MarshalEvent(ev Event) ([]byte, error) {
	var t string
	switch ev.(type) {
	case RevealCardEvent:
		t = "revealCard"
	case RoundScoreEvent:
		t = "roundScored"
	case RoundBeginsEvent:
		t = "roundBegins"
	default:
		return nil, fmt.Errorf("unknown event %T", ev)
	}
	data, err := json.Marshal(ev)
	if err != nil {
		return nil, err
	}
	return json.Marshal(envelope{Type: t, Data: data})
}

func UnmarshalEvent(b []byte) (Event, error) {
	var env envelope
	if err := json.Unmarshal(b, &env); err != nil {
		return nil, err
	}
	switch env.Type {
	case "revealCard":
		var e RevealCardEvent
		return e, json.Unmarshal(env.Data, &e)
	case "roundScored":
		var e RoundScoreEvent
		return e, json.Unmarshal(env.Data, &e)
	case "roundBegins":
		var e RoundBeginsEvent
		return e, json.Unmarshal(env.Data, &e)
	}
	return nil, fmt.Errorf("unknown event type %q", env.Type)
}

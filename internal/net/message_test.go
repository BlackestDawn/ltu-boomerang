package net

import (
	gonet "net"
	"reflect"
	"testing"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

var testCard = core.Card{
	Name:   "Uluru",
	Site:   'E',
	Region: "Northern Territory",
	Value:  4,
	Symbols: map[core.Category]*core.Symbol{
		core.Animal:   {Name: "Emus", Value: 0},
		core.Activity: {Name: "Indigenous Culture", Value: 0},
	},
}

func TestPayloadRoundTrip(t *testing.T) {
	cases := []struct {
		typ  MessageType
		in   any
		into any
	}{
		{Hello, HelloData{Version: ProtocolVersion}, &HelloData{}},
		{Welcome, WelcomeData{Seat: 2, NumPlayers: 4, Rounds: 4, Edition: "EU"}, &WelcomeData{}},
		{AskDraw, AskCardData{Hand: []core.Card{testCard}}, &AskCardData{}},
		{AskActivity, AskActivityData{
			Draft:   []core.Card{testCard},
			Options: []core.Symbol{{Name: "Swimming"}},
			Counts:  []int{2},
		}, &AskActivityData{}},
		{Answer, AnswerData{Choice: -1}, &AnswerData{}},
		{Final, FinalData{
			ScoreSheets: []core.ScoreSheet{{CategoryScores: map[core.Category][]int{core.Site: {3, 5}}}},
			Winner:      0,
		}, &FinalData{}},
		{Error, ErrorData{Message: "nope"}, &ErrorData{}},
	}

	a, b := gonet.Pipe()
	sender, receiver := NewTransport(a), NewTransport(b)
	defer sender.Close()
	defer receiver.Close()

	for _, c := range cases {
		msg, err := NewMessage(c.typ, c.in)
		if err != nil {
			t.Fatalf("%s: NewMessage: %v", c.typ, err)
		}
		go sender.Send(msg)
		got, err := receiver.Receive()
		if err != nil {
			t.Fatalf("%s: Receive: %v", c.typ, err)
		}
		if got.Type != c.typ {
			t.Fatalf("got type %s, want %s", got.Type, c.typ)
		}
		if err := got.Decode(c.into); err != nil {
			t.Fatalf("%s: Decode: %v", c.typ, err)
		}
		if out := reflect.ValueOf(c.into).Elem().Interface(); !reflect.DeepEqual(out, c.in) {
			t.Errorf("%s: got %+v, want %+v", c.typ, out, c.in)
		}
	}
}

func TestEventMessageRoundTrip(t *testing.T) {
	events := []core.Event{
		core.RoundBeginsEvent(3),
		core.RevealCardEvent{Player: 1, Card: testCard},
		core.RevealDraftEvent{Player: 2, Cards: []core.Card{testCard, testCard}},
		core.RoundScoreEvent{0: {core.Site: 4}, 1: {core.Animal: 2}},
	}
	for _, ev := range events {
		msg, err := NewEventMessage(ev)
		if err != nil {
			t.Fatalf("NewEventMessage(%T): %v", ev, err)
		}
		got, err := core.UnmarshalEvent(msg.Data)
		if err != nil {
			t.Fatalf("UnmarshalEvent(%T): %v", ev, err)
		}
		if !reflect.DeepEqual(got, ev) {
			t.Errorf("got %+v, want %+v", got, ev)
		}
	}
}

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition"
	"github.com/BlackestDawn/ltu-boomerang/internal/net"
	"github.com/BlackestDawn/ltu-boomerang/internal/scoring"
	"github.com/BlackestDawn/ltu-boomerang/internal/ui"
)

// runClient joins a game hosted at connectTo and plays it until the host
// sends the final results or the connection is lost.
func runClient(connectTo string) error {
	t, err := net.Dial(connectTo)
	if err != nil {
		return fmt.Errorf("could not connect to %s: %w", connectTo, err)
	}
	defer t.Close()

	hello, err := net.NewMessage(net.Hello, net.HelloData{Version: net.ProtocolVersion})
	if err != nil {
		return err
	}
	if err := t.Send(hello); err != nil {
		return fmt.Errorf("could not greet host: %w", err)
	}

	msg, err := t.Receive()
	if err != nil {
		return fmt.Errorf("host did not answer: %w", err)
	}
	if msg.Type == net.Error {
		return hostError(msg)
	}
	var welcome net.WelcomeData
	if msg.Type != net.Welcome || msg.Decode(&welcome) != nil {
		return fmt.Errorf("unexpected %s from host", msg.Type)
	}
	// Loading the edition sets the category names shown on score sheets.
	if _, err := edition.GetEdition(welcome.Edition); err != nil {
		return fmt.Errorf("host plays an unsupported edition: %w", err)
	}

	textUI := ui.NewTextUI(os.Stdout, os.Stdin)
	textUI.SetSelf(welcome.Seat)
	textUI.SetNumPlayers(welcome.NumPlayers)
	fmt.Printf("Connected as player %d, waiting for the game to start\n", welcome.Seat+1)

	c := &client{t: t, ui: textUI}
	return c.play()
}

type client struct {
	t  net.Transport
	ui *ui.TextUI
}

// play handles messages from the host one at a time. Questions are answered
// before the next message is read, so stdin is only read here.
func (c *client) play() error {
	for {
		msg, err := c.t.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return errors.New("host closed the connection")
			}
			return fmt.Errorf("connection to host lost: %w", err)
		}

		switch msg.Type {
		case net.Event:
			ev, err := core.UnmarshalEvent(msg.Data)
			if err != nil {
				return fmt.Errorf("bad event from host: %w", err)
			}
			c.ui.OnEvent(ev)

		case net.AskThrow, net.AskDraw:
			var ask net.AskCardData
			if err := msg.Decode(&ask); err != nil {
				return fmt.Errorf("bad question from host: %w", err)
			}
			question := "Choose your throw card"
			if msg.Type == net.AskDraw {
				question = "Choose a card to draft"
			}
			if err := c.answer(c.ui.AskForCard(ask.Hand, question)); err != nil {
				return err
			}

		case net.AskActivity:
			var ask net.AskActivityData
			if err := msg.Decode(&ask); err != nil {
				return fmt.Errorf("bad question from host: %w", err)
			}
			if err := c.answer(c.ui.AskForActivity(ask.Draft, ask.Options, ask.Counts)); err != nil {
				return err
			}

		case net.Final:
			var final net.FinalData
			if err := msg.Decode(&final); err != nil {
				return fmt.Errorf("bad results from host: %w", err)
			}
			c.ui.ShowFinalResults(finalResult(final))

		case net.Error:
			// The host repeats the question after an error, so keep going.
			fmt.Println(hostError(msg))

		case net.Bye:
			return nil
		}
	}
}

func (c *client) answer(choice int) error {
	msg, err := net.NewMessage(net.Answer, net.AnswerData{Choice: choice})
	if err != nil {
		return err
	}
	if err := c.t.Send(msg); err != nil {
		return fmt.Errorf("could not send answer: %w", err)
	}
	return nil
}

// finalResult rebuilds the session result the UI expects. Only the score
// sheets are sent, which is all the results screen uses.
func finalResult(f net.FinalData) scoring.SessionResult {
	scores := make([]*core.ScoreSheet, len(f.ScoreSheets))
	for i := range f.ScoreSheets {
		scores[i] = &f.ScoreSheets[i]
	}
	return scoring.SessionResult{Scores: scores, Winner: f.Winner}
}

func hostError(msg net.Message) error {
	var e net.ErrorData
	if err := msg.Decode(&e); err != nil {
		return errors.New("host sent an unreadable error")
	}
	return fmt.Errorf("host: %s", e.Message)
}

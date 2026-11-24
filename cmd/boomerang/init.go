package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/BlackestDawn/ltu-boomerang/internal/agent"
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/engine"
	"github.com/BlackestDawn/ltu-boomerang/internal/net"
	"github.com/BlackestDawn/ltu-boomerang/internal/ui"
)

// helloTimeout is how long a new connection has to introduce itself.
const helloTimeout = 10 * time.Second

// buildControllers seats the local player first, then the remote players in
// the order they connect, then the bots. The remote controllers are also
// returned on their own so they can be finished after the game.
func buildControllers(cfg *Config, ui *ui.TextUI, rng *rand.Rand, es core.EventServer) ([]engine.PlayerController, []*agent.RemoteController, error) {
	numRemote := cfg.PlayerCount - 1
	numPlayers := cfg.PlayerCount + cfg.BotCount

	participants := []engine.PlayerController{agent.NewHumanController(ui, es.Subscribe())}

	remotes, err := acceptRemotes(":"+cfg.Listen, numRemote, numPlayers, cfg.Rounds, cfg.Edition, rng, es)
	if err != nil {
		return nil, nil, err
	}
	for _, r := range remotes {
		participants = append(participants, r)
	}

	for i := 0; i < cfg.BotCount; i++ {
		participants = append(participants, agent.NewBotController(rng))
	}
	return participants, remotes, nil
}

// acceptRemotes waits on addr until count clients have connected and been
// welcomed. Connections that fail the handshake are dropped and their slot
// stays open.
func acceptRemotes(addr string, count, numPlayers, rounds int, edition string, rng *rand.Rand, es core.EventServer) ([]*agent.RemoteController, error) {
	if count <= 0 {
		return nil, nil
	}
	ln, err := net.Listen(addr)
	if err != nil {
		return nil, err
	}
	defer ln.Close()

	var remotes []*agent.RemoteController
	for len(remotes) < count {
		fmt.Printf("Waiting for remote players on %s (%d/%d connected)\n", ln.Addr(), len(remotes), count)
		t, err := ln.Accept()
		if err != nil {
			return nil, fmt.Errorf("accepting remote player: %w", err)
		}
		seat := len(remotes) + 1
		if err := handshake(t, seat, numPlayers, rounds, edition); err != nil {
			fmt.Printf("Dropped connection: %s\n", err)
			t.Close()
			continue
		}
		remotes = append(remotes, agent.NewRemoteController(t, es, agent.NewBotController(rng)))
	}
	return remotes, nil
}

// handshake waits for the client's Hello and answers with a Welcome giving it
// its seat and the edition being played.
func handshake(t net.Transport, seat, numPlayers, rounds int, edition string) error {
	t.SetReadDeadline(time.Now().Add(helloTimeout))
	defer t.SetReadDeadline(time.Time{})

	msg, err := t.Receive()
	if err != nil {
		return fmt.Errorf("no hello: %w", err)
	}
	var hello net.HelloData
	if msg.Type != net.Hello || msg.Decode(&hello) != nil {
		return fmt.Errorf("expected hello, got %s", msg.Type)
	}
	if hello.Version != net.ProtocolVersion {
		refusal, _ := net.NewMessage(net.Error, net.ErrorData{
			Message: fmt.Sprintf("protocol version %d not supported, host runs %d", hello.Version, net.ProtocolVersion),
		})
		t.Send(refusal)
		return fmt.Errorf("client has protocol version %d", hello.Version)
	}

	welcome, err := net.NewMessage(net.Welcome, net.WelcomeData{
		Seat:       seat,
		NumPlayers: numPlayers,
		Rounds:     rounds,
		Edition:    edition,
	})
	if err != nil {
		return err
	}
	return t.Send(welcome)
}

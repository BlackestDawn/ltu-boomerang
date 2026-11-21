package main

import (
	"fmt"
	"math/rand"

	"github.com/BlackestDawn/ltu-boomerang/internal/agent"
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/engine"
	"github.com/BlackestDawn/ltu-boomerang/internal/ui"
)

func buildControllers(cfg *Config, ui *ui.TextUI, rng *rand.Rand, es core.EventServer) []engine.PlayerController {
	var participants []engine.PlayerController
	participants = append(participants, agent.NewHumanController(ui, es.Subscribe()))
	if cfg.BotCount != nil {
		for i := 0; i < *cfg.BotCount; i++ {
			participants = append(participants, agent.NewBotController(rng))
		}
	}
	if *cfg.PlayerCount > 1 {
		fmt.Println("Waiting for remote players to connect")
		// TODO: add remote players
	}
	return participants
}

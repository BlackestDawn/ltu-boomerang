package main

import (
	"math/rand"

	"github.com/BlackestDawn/ltu-boomerang/internal/agent"
	"github.com/BlackestDawn/ltu-boomerang/internal/engine"
	"github.com/BlackestDawn/ltu-boomerang/internal/ui"
)

func buildControllers(cfg *Config, ui *ui.TextUI, rng *rand.Rand) []engine.PlayerController {
	var participants []engine.PlayerController
	participants = append(participants, agent.NewHumanController(ui))
	for i := 0; i < *cfg.BotCount; i++ {
		participants = append(participants, agent.NewBotController(rng))
	}
	// TODO: add remote players
	return participants
}

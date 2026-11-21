package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition"
	"github.com/BlackestDawn/ltu-boomerang/internal/engine"
	"github.com/BlackestDawn/ltu-boomerang/internal/ui"
)

func main() {
	cfg, err := buildConfig(os.Args)
	if err != nil {
		fmt.Printf("Could not initialize config: %s\n", err)
		os.Exit(1)
	}

	if *cfg.ConnectTo != "" {
		runClient(*cfg.ConnectTo)
		return
	}

	ed, err := edition.GetEdition(*cfg.Edition)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ui := ui.NewTextUI(os.Stdout, os.Stdin)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventBus := make(chan core.Event, 16)
	es := core.NewEventServer(ctx, eventBus)

	players := buildControllers(cfg, ui, rng, es)
	ui.SetNumPlayers(len(players))
	g := engine.NewGame(ed, players, rng, *cfg.PassingDir, *cfg.Rounds, eventBus)

	result := g.Play()

	ui.ShowFinalResults(result)

	cfg.cleanup()
	fmt.Println("Thank you for playing, come again.")
}

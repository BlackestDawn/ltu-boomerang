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
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if *cfg.ConnectTo != "" {
		if err := runClient(*cfg.ConnectTo); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Thank you for playing, come again.")
		return
	}

	ed, err := edition.GetEdition(*cfg.Edition)
	if err != nil {
		fmt.Printf("Could not load edition: %s\n", err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ui := ui.NewTextUI(os.Stdout, os.Stdin)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	es := core.NewEventServer(ctx)

	players, remotes, err := buildControllers(cfg, ui, rng, es)
	if err != nil {
		fmt.Printf("Could not set up players: %s\n", err)
		os.Exit(1)
	}
	ui.SetNumPlayers(len(players))
	g := engine.NewGame(ed, players, rng, *cfg.PassingDir, *cfg.Rounds, es)

	result := g.Play()

	for _, r := range remotes {
		r.Finish(result)
	}
	ui.ShowFinalResults(result)

	fmt.Println("Thank you for playing, come again.")
}

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

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
	ui := &ui.TextUI{}
	players := buildControllers(cfg, ui, rng)
	g := engine.NewGame(ed, players, rng, *cfg.PassingDir)

	result := g.Play()

	ui.ShowResults(result)

	cfg.cleanup()
	fmt.Println("Thank you for playing, come again.")
}

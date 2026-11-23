package main

import (
	"fmt"
	"strconv"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/akamensky/argparse"
)

type Config struct {
	ConnectTo   *string
	Listen      *string
	Edition     *string
	PassingDir  *string
	PlayerCount *int
	BotCount    *int
	Rounds      *int
}

func buildConfig(args []string) (*Config, error) {
	parser := argparse.NewParser("boomerang", "Boomerang Cardgame")

	edition := parser.Selector(
		"e",
		"edition",
		[]string{"AU", "EU", "US"},
		&argparse.Options{
			Help:    "Select which edition to use",
			Default: "AU",
		},
	)
	passing := parser.Selector(
		"d",
		"direction",
		[]string{"left", "alternate", "right"},
		&argparse.Options{
			Help:    "Direction to pass cards during each round",
			Default: "left",
		},
	)
	players := parser.Int(
		"p",
		"players",
		&argparse.Options{
			Help:     fmt.Sprintf("Number of players: 1-%d", core.MaxPlayerCount),
			Default:  2,
			Validate: func(args []string) error { return validateCount(args[0], core.MaxPlayerCount, core.MinPlayerCount) },
		},
	)
	bots := parser.Int(
		"b",
		"bots",
		&argparse.Options{
			Help:     fmt.Sprintf("Number of bots, can only fill remaining \"player\" slots, max would become %d - player count", core.MaxPlayerCount),
			Default:  0,
			Validate: func(args []string) error { return validateCount(args[0], core.MaxPlayerCount-*players, 0) },
		},
	)
	connectTo := parser.String(
		"c",
		"connect",
		&argparse.Options{
			Help: "Connect to a remote session. If specified no other option will have any effect",
		},
	)
	listen := parser.String(
		"l",
		"listen",
		&argparse.Options{
			Help:    "Address to listen for remote players on, used when there are more players than 1",
			Default: core.ListenPort,
		},
	)
	numRounds := parser.Int(
		"r",
		"rounds",
		&argparse.Options{
			Help:     fmt.Sprintf("Number of rounds to play: %d-%d", core.MinRounds, core.MaxRounds),
			Default:  core.MinRounds,
			Validate: func(args []string) error { return validateCount(args[0], core.MaxRounds, core.MinRounds) },
		},
	)

	err := parser.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("%s", parser.Usage(err))
	}

	return &Config{
		ConnectTo:   connectTo,
		Listen:      listen,
		Edition:     edition,
		PassingDir:  passing,
		PlayerCount: players,
		BotCount:    bots,
		Rounds:      numRounds,
	}, nil
}

func validateCount(input string, max int, min int) error {
	count, err := strconv.Atoi(input)
	if err != nil {
		return err
	}

	if count < min || count > max {
		return fmt.Errorf("Invalid value: %d\n", count)
	}

	return nil
}

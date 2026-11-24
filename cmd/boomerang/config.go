package main

import (
	"fmt"
	"strconv"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/akamensky/argparse"
)

type Config struct {
	ConnectTo   string
	Listen      string
	Edition     string
	PassingDir  string
	PlayerCount int
	BotCount    int
	Rounds      int
}

func buildConfig(args []string) (*Config, error) {
	parser := argparse.NewParser(
		"boomerang",
		"Boomerang: a draft-and-write card game for 2-4 players (Australia, Europe and USA editions).\n"+
			"Run without --connect to host a game: you play as player 1, remote players join over TCP\n"+
			"and bots fill any empty seats. Run with --connect <address> to join a game someone else is hosting.",
	)

	edition := parser.Selector(
		"e",
		"edition",
		[]string{"AU", "EU", "US"},
		&argparse.Options{
			Help:    "Edition to play: AU (Australia), EU (Europe), US (USA)",
			Default: "AU",
		},
	)
	passing := parser.Selector(
		"d",
		"direction",
		[]string{"left", "alternate", "right"},
		&argparse.Options{
			Help:    "Direction hands are passed: left (next player), right (previous player), alternate (switches each round)",
			Default: "left",
		},
	)
	players := parser.Int(
		"p",
		"players",
		&argparse.Options{
			Help:     fmt.Sprintf("Number of human players including you, %d-%d. Everyone except you joins with --connect", core.MinPlayerCount, core.MaxPlayerCount),
			Default:  2,
			Validate: func(args []string) error { return validateCount(args[0], core.MinPlayerCount, core.MaxPlayerCount) },
		},
	)
	bots := parser.Int(
		"b",
		"bots",
		&argparse.Options{
			Help:     fmt.Sprintf("Number of bots to fill empty seats, up to %d. At least 1 is needed when playing alone", core.MaxPlayerCount-1),
			Default:  0,
			Validate: func(args []string) error { return validateCount(args[0], 0, core.MaxPlayerCount-*players) },
		},
	)
	connectTo := parser.String(
		"c",
		"connect",
		&argparse.Options{
			Help: fmt.Sprintf("Join a game at <address>[:port] instead of hosting one (port defaults to %s). All other options are ignored", core.ListenPort),
		},
	)
	listen := parser.String(
		"l",
		"listen",
		&argparse.Options{
			Help:    "Port to wait for remote players on. Only used when there are more than 1 player",
			Default: core.ListenPort,
		},
	)
	numRounds := parser.Int(
		"r",
		"rounds",
		&argparse.Options{
			Help:     fmt.Sprintf("Number of rounds to play, %d-%d. The standard rules use %d", core.MinRounds, core.MaxRounds, core.MinRounds),
			Default:  core.MinRounds,
			Validate: func(args []string) error { return validateCount(args[0], core.MinRounds, core.MaxRounds) },
		},
	)

	err := parser.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("%s", parser.Usage(err.Error()+"\n"))
	}

	if *players+*bots < 2 {
		return nil, fmt.Errorf("%s", parser.Usage("total between players and bots must be at least 2\n"))
	}

	return &Config{
		ConnectTo:   *connectTo,
		Listen:      *listen,
		Edition:     *edition,
		PassingDir:  *passing,
		PlayerCount: *players,
		BotCount:    *bots,
		Rounds:      *numRounds,
	}, nil
}

func validateCount(input string, min, max int) error {
	count, err := strconv.Atoi(input)
	if err != nil {
		return err
	}

	if count < min || count > max {
		return fmt.Errorf("invalid value: %d", count)
	}

	return nil
}

package main

import (
	"fmt"
	"strconv"

	"github.com/akamensky/argparse"
)

type Config struct {
	ConnectTo *string
	// Ruleset     *string
	PassingDir  *string
	PlayerCount *int
	BotCount    *int
}

func (c *Config) cleanup() {
	// TODO: cleanup stuff like networking
}

func buildConfig(args []string) (*Config, error) {
	parser := argparse.NewParser("boomerang", "Boomerang Cardgame")

	// rules := parser.Selector("r", "rules", []string{"AU", "EU", "US"}, &argparse.Options{Help: "Select which ruleset to use", Default: "AU"})
	passing := parser.Selector("d", "direction", []string{"left", "alternate", "right"}, &argparse.Options{Help: "Direction to pass cards during each round", Default: "left"})
	players := parser.Int("p", "players", &argparse.Options{Help: "Number of players: 1-4", Default: 2, Validate: func(args []string) error { return validateCount(args[0], maxPlayerCount, minPlayerCount) }})
	bots := parser.Int("b", "bots", &argparse.Options{Help: "Number of bots, can only fill remaining \"player\" slots", Default: 0, Validate: func(args []string) error { return validateCount(args[0], maxPlayerCount-*players, 0) }})
	connectTo := parser.String("c", "connect", &argparse.Options{Help: "Connect to a remote session. If specified no other option will have any effect"})

	err := parser.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("%s", parser.Usage(err))
	}

	return &Config{
		ConnectTo: connectTo,
		// Ruleset:     rules,
		PassingDir:  passing,
		PlayerCount: players,
		BotCount:    bots,
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

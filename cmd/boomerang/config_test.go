package main

import (
	"strings"
	"testing"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

func parse(args string) (*Config, error) {
	return buildConfig(append([]string{"boomerang"}, strings.Fields(args)...))
}

// Requirement 1: a game has between 2 and 4 players in total (humans + bots).
func TestPlayerCountLimits(t *testing.T) {
	cases := []struct {
		args string
		ok   bool
	}{
		{"", true}, // default: 2 players
		{"-p 1 -b 1", true},
		{"-p 2 -b 2", true},
		{"-p 4", true},
		{"-p 1 -b 3", true},
		{"-p 1 -b 0", false}, // only one player
		{"-p 0 -b 2", false}, // nobody at the host
		{"-p 5", false},
		{"-p 3 -b 2", false}, // five in total
		{"-p 1 -b 4", false},
		{"-b -1", false},
	}
	for _, c := range cases {
		cfg, err := parse(c.args)
		if c.ok && err != nil {
			t.Errorf("%q rejected: %v", c.args, err)
			continue
		}
		if !c.ok {
			if err == nil {
				t.Errorf("%q accepted with %d players and %d bots, want error",
					c.args, *cfg.PlayerCount, *cfg.BotCount)
			}
			continue
		}
		if total := *cfg.PlayerCount + *cfg.BotCount; total < 2 || total > core.MaxPlayerCount {
			t.Errorf("%q gives %d players in total", c.args, total)
		}
	}
}

// Requirement 12: the game is four rounds by default.
func TestDefaultsToFourRounds(t *testing.T) {
	cfg, err := parse("")
	if err != nil {
		t.Fatal(err)
	}
	if *cfg.Rounds != 4 {
		t.Errorf("default rounds %d, want 4", *cfg.Rounds)
	}
	if *cfg.Edition != "AU" || *cfg.PassingDir != "left" {
		t.Errorf("defaults: edition %q passing %q, want AU/left", *cfg.Edition, *cfg.PassingDir)
	}
}

func TestValidateCount(t *testing.T) {
	cases := []struct {
		in string
		ok bool
	}{
		{"2", true}, {"4", true}, {"1", false}, {"5", false}, {"x", false},
	}
	for _, c := range cases {
		if err := validateCount(c.in, 4, 2); (err == nil) != c.ok {
			t.Errorf("validateCount(%q) error = %v, want ok=%v", c.in, err, c.ok)
		}
	}
}

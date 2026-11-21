package ui

import (
	"bufio"
	"io"
	"sync"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

type TextUI struct {
	out io.Writer
	in  *bufio.Reader

	mu         sync.Mutex // guards the fields below
	draftCards map[int][]core.Card
	numPlayers int
	round      int         // current round, 0-based
	totals     map[int]int // running score per player
	pending    block       // round result screen not yet shown to the player
	finished   bool        // final results shown, ignore late round results
}

func NewTextUI(out io.Writer, in io.Reader) *TextUI {
	return &TextUI{
		out:        out,
		in:         bufio.NewReader(in),
		draftCards: make(map[int][]core.Card),
		totals:     make(map[int]int),
	}
}

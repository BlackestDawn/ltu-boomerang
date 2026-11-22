package ui

import (
	"fmt"
	"strings"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

// seatOf places players around the table clockwise from the local player, so
// passing to the next player hands cards to the left.
func seatOf(player, numPlayers int) seat {
	rel := (player - localPlayer + numPlayers) % numPlayers
	switch numPlayers {
	case 2:
		return []seat{seatBottom, seatTop}[rel]
	case 3:
		return []seat{seatBottom, seatLeft, seatRight}[rel]
	default:
		return seat(rel % 4)
	}
}

// playerName is how a player is referred to on screen.
func playerName(p int) string {
	if p == localPlayer {
		return "You"
	}
	return fmt.Sprintf("Player %d", p+1)
}

// PrintTable clears the screen and prints the banner and the table with
// every player's draft along its edges and nothing in the middle.
func (ui *TextUI) PrintTable() {
	ui.printScreen(nil)
}

// PrintHand prints the cards in a grid.
func (ui *TextUI) PrintHand(cards []core.Card) {
	fmt.Fprintln(ui.out)
	if len(cards) > 0 {
		printCardsGrid(ui.out, cards, 0, handColumns, "Your hand")
	} else {
		fmt.Fprintln(ui.out, "    No cards in hand")
	}
	fmt.Fprintln(ui.out)
}

// printScreen clears the screen and prints the banner and the table with
// middle in its centre, see printLines for cursor handling.
func (ui *TextUI) printScreen(middle block) (done func()) {
	ui.mu.Lock()
	screen := ui.composeScreen(middle)
	ui.mu.Unlock()
	return ui.printLines(screen)
}

// composeScreen puts the banner above the table with middle in its centre.
// The caller must hold ui.mu.
func (ui *TextUI) composeScreen(middle block) block {
	return vjoin(1, block(strings.Split(banner, "\n")), ui.renderTable(middle))
}

// printLines clears the screen and prints lines. If they contain cursorMark,
// the cursor is left there and the returned function moves it back below the
// lines; call it once input has been read. Otherwise it does nothing.
func (ui *TextUI) printLines(screen block) (done func()) {
	lines, row, col, ok := findCursor(screen)

	fmt.Fprint(ui.out, clearScreen)
	for _, l := range lines {
		fmt.Fprintln(ui.out, strings.TrimRight(l, " "))
	}
	if !ok {
		return func() {}
	}

	// The cursor is at the start of the line below the table.
	fmt.Fprintf(ui.out, "\033[%dA\033[%dG", len(lines)-row, col+1)
	return func() {
		// Entering the input moved the cursor to the line below the prompt.
		if down := len(lines) - row - 1; down > 0 {
			fmt.Fprintf(ui.out, "\033[%dB", down)
		}
		fmt.Fprint(ui.out, "\r")
	}
}

// renderTable draws the drafts of all players along the edges of a square
// table, the local player at the bottom, with middle in the centre.
// The caller must hold ui.mu.
func (ui *TextUI) renderTable(middle block) block {
	numPlayers := ui.numPlayers
	for p := range ui.draftCards {
		numPlayers = max(numPlayers, p+1)
	}
	numPlayers = max(numPlayers, 1)

	edges := make(map[seat]block)
	for p := range numPlayers {
		header := playerName(p)
		if p == localPlayer {
			header = "Your draft"
		}
		s := seatOf(p, numPlayers)
		switch s {
		case seatTop, seatBottom:
			edges[s] = renderGrid(ui.draftCards[p], -1, 0, header)
		default:
			edges[s] = renderGrid(ui.draftCards[p], 0, sideColumns, header)
		}
	}

	left, right := edges[seatLeft], edges[seatRight]
	// Give both sides the same width so the middle stays centred, even when
	// only one side is occupied.
	if len(left) > 0 || len(right) > 0 {
		w := max(left.width(), right.width())
		left = vjoin(0, left, blankBlock(w))
		right = vjoin(0, right, blankBlock(w))
	}
	middle = vjoin(0, middle, blankBlock(16))

	return vjoin(1, edges[seatTop], hjoin(tableGap, left, middle, right), edges[seatBottom])
}

// blankBlock is a single empty line of the given width.
func blankBlock(width int) block {
	return block{strings.Repeat(" ", width)}
}

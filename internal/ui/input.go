package ui

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

func (ui *TextUI) AskForNumber(max, min int) int {
	for {
		fmt.Fprintf(ui.out, "Choose a number between %d and %d: ", min, max)
		choice, err := ui.readNumber(min, max)
		if err == nil || errors.Is(err, io.EOF) {
			return choice
		}
		fmt.Fprintln(ui.out, err)
	}
}

// AskForCard shows the table with the cards in the middle and asks the player
// to pick one of them. It returns the index of the chosen card.
func (ui *TextUI) AskForCard(cards []core.Card, question string) int {
	labels := make([]string, len(cards))
	for i := range cards {
		labels[i] = fmt.Sprintf("[%d]", i+1)
	}
	grid := block(renderLabeledGrid(cards, 0, handColumns, "Your hand", labels))
	prompt := fmt.Sprintf("%s (1-%d): ", question, len(cards))

	return ui.askInMiddle(1, len(cards), func(msg string) block {
		return vjoin(1, grid, block{msg}, block{prompt + cursorMark})
	}) - 1
}

// AskForActivity shows the table with the draft and the given activity
// options in the middle and asks the player to pick one or skip. counts[i] is
// how many of options[i] the draft holds. It returns the index of the chosen
// option, or -1 when skipping.
func (ui *TextUI) AskForActivity(draft []core.Card, options []core.Symbol, counts []int) int {
	grid := block(renderGrid(draft, 0, handColumns, "Your draft"))
	list := block{rule("Activities", 30)}
	for i, o := range options {
		list = append(list, fmt.Sprintf("%s%d)%s %-16s ×%d", ansiBlue, i+1, ansiReset, o.Name, counts[i]))
	}
	list = append(list, fmt.Sprintf("%s0)%s skip", ansiBlue, ansiReset))
	prompt := fmt.Sprintf("Score which activity (0-%d): ", len(options))

	return ui.askInMiddle(0, len(options), func(msg string) block {
		return vjoin(1, grid, list, block{msg}, block{prompt + cursorMark})
	}) - 1
}

// askInMiddle prints the table with the block built by middle in its centre,
// puts the cursor on the mark in it and reads a number in [min, max]. On bad
// input the screen is redrawn with the error passed to middle.
func (ui *TextUI) askInMiddle(min, max int, middle func(msg string) block) int {
	ui.flushPending()
	msg := ""
	for {
		done := ui.printScreen(middle(ansiYellow + msg + ansiReset))
		choice, err := ui.readNumber(min, max)
		done()
		if err == nil || errors.Is(err, io.EOF) {
			return choice
		}
		msg = err.Error()
	}
}

// flushPending shows a round result that arrived since the last prompt and
// waits for the player to press Enter. Runs on the goroutine reading input, so
// stdin is never read from two places at once.
func (ui *TextUI) flushPending() {
	ui.mu.Lock()
	screen := ui.pending
	ui.pending = nil
	ui.mu.Unlock()
	if screen == nil {
		return
	}
	done := ui.printLines(screen)
	ui.readLine()
	done()
}

// readLine reads one line of input without the line ending.
func (ui *TextUI) readLine() (string, error) {
	line, err := ui.in.ReadString('\n')
	if errors.Is(err, io.EOF) && line != "" {
		err = nil
	}
	return strings.TrimSpace(line), err
}

// readNumber reads a line holding a number in [min, max]. On io.EOF it
// returns min, so a closed input doesn't make callers loop forever.
func (ui *TextUI) readNumber(min, max int) (int, error) {
	input, err := ui.readLine()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return min, err
		}
		return 0, fmt.Errorf("Could not properly read input: %w", err)
	}
	choice, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("Not a number: %q", input)
	}
	if choice < min || choice > max {
		return 0, fmt.Errorf("Input outside of range: %d", choice)
	}
	return choice, nil
}

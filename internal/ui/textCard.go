package ui

import (
	"fmt"
	"io"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

// Symbol printing order and their coloration
var symbolOrder = []core.Category{core.Collection, core.Animal, core.Activity}

var symbolColor = map[core.Category]string{
	core.Collection: ansiGreen,
	core.Animal:     ansiYellow,
	core.Activity:   ansiBlue,
}

// renderCard returns the card as lines of text, all the same visible width.
func renderCard(c core.Card) []string {
	border := strings.Repeat("─", cardInnerWidth)
	lines := []string{"╭" + border + "╮"}

	// Top corners: card value on the left, site on the right.
	lines = append(lines, cardLine(spread(fmt.Sprint(c.Value), siteLabel(c.Site)), ansiBold))
	// lines = append(lines, cardLine("", ""))
	lines = append(lines, cardLine(center(c.Name), ansiBold))
	lines = append(lines, cardLine(center(c.Region), ""))
	lines = append(lines, cardLine("", ""))

	// Always reserve two symbol lines
	drawn := 0
	for _, cat := range symbolOrder {
		s, ok := c.Symbols[cat]
		if !ok || s == nil || drawn == 2 {
			continue
		}
		lines = append(lines, cardLine(spread("● "+s.Name, fmt.Sprint(s.Value)), symbolColor[cat]))
		drawn++
	}
	for ; drawn < 2; drawn++ {
		lines = append(lines, cardLine("", ""))
	}

	// Bottom corners mirror the top, like a real playing card.
	lines = append(lines, cardLine(spread(siteLabel(c.Site), fmt.Sprint(c.Value)), ansiBold))
	lines = append(lines, "╰"+border+"╯")
	return lines
}

// printCardsGrid prints the cards as a grid, see renderGrid for the layout rules.
func printCardsGrid(w io.Writer, cards []core.Card, rows, cols int, header string) {
	for _, line := range renderGrid(cards, rows, cols, header) {
		fmt.Fprintln(w, line)
	}
}

// renderGrid lays the cards out in a grid and returns it as lines of equal
// visible width. rows and cols each take:
//
//	-1  a single row/column (if both are -1, a single row wins)
//	 0  grow as needed to fit all cards
//	>0  that many rows/columns
//
// If both are fixed and too small, rows are added. The grid never gets more
// rows or columns than the cards need. A non-empty header is drawn as a rule
// above the cards, e.g. to show which player they belong to.
func renderGrid(cards []core.Card, rows, cols int, header string) []string {
	return renderLabeledGrid(cards, rows, cols, header, nil)
}

// renderLabeledGrid is renderGrid with an optional label centred above each card.
func renderLabeledGrid(cards []core.Card, rows, cols int, header string, labels []string) []string {
	rows, cols = gridShape(len(cards), rows, cols)

	cardWidth := cardInnerWidth + 2
	width := cols*cardWidth + max(cols-1, 0)*len(cardGap)
	if header != "" {
		width = max(width, utf8.RuneCountInString(header)+6)
	}

	var lines []string
	if header != "" {
		lines = append(lines, rule(header, width))
	}
	if len(cards) == 0 {
		return append(lines, centerIn("(no cards)", width))
	}

	for r := range rows {
		from := r * cols
		to := min(from+cols, len(cards))
		rendered := make([][]string, 0, to-from)
		for _, c := range cards[from:to] {
			rendered = append(rendered, renderCard(c))
		}
		if labels != nil {
			parts := make([]string, 0, to-from)
			for i := from; i < to; i++ {
				label := ""
				if i < len(labels) {
					label = labels[i]
				}
				parts = append(parts, ansiBold+centerIn(label, cardWidth)+ansiReset)
			}
			lines = append(lines, padTo(strings.Join(parts, cardGap), width))
		}
		for line := range rendered[0] {
			parts := make([]string, len(rendered))
			for i := range rendered {
				parts[i] = rendered[i][line]
			}
			lines = append(lines, padTo(strings.Join(parts, cardGap), width))
		}
	}
	return lines
}

// gridShape resolves the rows/cols arguments of renderGrid for n cards.
func gridShape(n, rows, cols int) (int, int) {
	if n == 0 {
		return 0, 0
	}
	switch {
	case rows == -1:
		rows, cols = 1, n
	case cols == -1:
		rows, cols = n, 1
	case rows <= 0 && cols <= 0:
		cols = int(math.Ceil(math.Sqrt(float64(n))))
	case cols <= 0:
		cols = ceilDiv(n, rows)
	}
	cols = min(cols, n)
	return ceilDiv(n, cols), cols
}

func ceilDiv(a, b int) int {
	return (a + b - 1) / b
}

// rule draws a horizontal line of the given width with text centred in it.
func rule(text string, width int) string {
	text = " " + text + " "
	free := max(width-utf8.RuneCountInString(text), 0)
	left := free / 2
	return ansiBold + strings.Repeat("─", left) + text + strings.Repeat("─", free-left) + ansiReset
}

// cardLine pads text to the inner width, colours it and adds the side borders.
// Padding happens before colouring so escape codes don't break the alignment.
func cardLine(text, color string) string {
	text = " " + fit(text, cardInnerWidth-2) + " "
	if color != "" {
		text = color + text + ansiReset
	}
	return "│" + text + "│"
}

// fit truncates or right-pads s to exactly n characters.
func fit(s string, n int) string {
	l := utf8.RuneCountInString(s)
	if l > n {
		r := []rune(s)
		return string(r[:n-1]) + "…"
	}
	return s + strings.Repeat(" ", n-l)
}

// spread puts left and right at opposite edges of the content area.
func spread(left, right string) string {
	space := cardInnerWidth - 2 - utf8.RuneCountInString(right)
	return fit(left, space-1) + " " + right
}

// center centres s in the content area.
func center(s string) string {
	n := cardInnerWidth - 2
	l := utf8.RuneCountInString(s)
	if l >= n {
		return s
	}
	return strings.Repeat(" ", (n-l)/2) + s
}

// siteLabel shows the site as a character when printable, otherwise as a number.
func siteLabel(b byte) string {
	if b >= '!' && b <= '~' {
		return string(b)
	}
	return fmt.Sprint(b)
}

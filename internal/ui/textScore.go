package ui

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/scoring"
)

// scoreOrder is the order categories are listed in on score sheets.
var scoreOrder = []core.Category{
	core.ThrowCatch, core.Site, core.RegionBonus,
	core.Collection, core.Animal, core.Activity,
	core.TravelBonus,
}

// ShowRoundResult prepares a screen with the round's scores in the middle of
// the table, drafts still showing. It is shown, waiting for Enter, before the
// player is next asked for input.
func (ui *TextUI) ShowRoundResult(result scoring.RoundScore) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	if ui.finished {
		return
	}

	players := sortedKeys(result)
	header := []string{fmt.Sprintf("Round %d", ui.round+1)}
	for _, p := range players {
		header = append(header, playerName(p, ui.self))
	}

	var catRows [][]int
	var labels []string
	for _, cat := range presentCategories(result) {
		row := make([]int, len(players))
		for i, p := range players {
			row[i] = result[p][cat]
		}
		labels = append(labels, cat.String())
		catRows = append(catRows, row)
	}

	roundRow := make([]int, len(players))
	totalRow := make([]int, len(players))
	for i, p := range players {
		for _, v := range result[p] {
			roundRow[i] += v
		}
		ui.totals[p] += roundRow[i]
		totalRow[i] = ui.totals[p]
	}

	table := renderScoreTable(header,
		scoreSection(labels, catRows),
		scoreSection([]string{"Round", "Total"}, [][]int{roundRow, totalRow}),
	)
	middle := vjoin(1, table, block{"Press Enter to continue " + cursorMark})
	ui.pending = ui.composeScreen(middle)
}

// ShowFinalResults clears the screen and shows the score per category and
// per round for every player, followed by the ranking and the winner.
func (ui *TextUI) ShowFinalResults(result scoring.SessionResult) {
	ui.mu.Lock()
	// Everything in the last round result is in the final score sheets too.
	ui.finished = true
	ui.pending = nil
	ui.mu.Unlock()

	scoreSheets := result.Scores
	header := []string{"Final score"}
	totals := make([]int, len(scoreSheets))
	for p, s := range scoreSheets {
		header = append(header, playerName(p, ui.self))
		totals[p] = s.Total()
	}

	// Category totals over the whole game.
	sheets := make(map[int]map[core.Category]int)
	numRounds := 0
	for p, s := range scoreSheets {
		sheets[p] = make(map[core.Category]int)
		for cat, scores := range s.CategoryScores {
			sheets[p][cat] = s.CategoryTotal(cat)
			numRounds = max(numRounds, len(scores))
		}
	}
	var catLabels []string
	var catRows [][]int
	for _, cat := range presentCategories(sheets) {
		row := make([]int, len(scoreSheets))
		for p := range scoreSheets {
			row[p] = sheets[p][cat]
		}
		catLabels = append(catLabels, cat.String())
		catRows = append(catRows, row)
	}

	// Score per round. Categories are added once per round, in round order.
	var roundLabels []string
	var roundRows [][]int
	for r := range numRounds {
		row := make([]int, len(scoreSheets))
		for p, s := range scoreSheets {
			for _, scores := range s.CategoryScores {
				if r < len(scores) {
					row[p] += scores[r]
				}
			}
		}
		roundLabels = append(roundLabels, fmt.Sprintf("Round %d", r+1))
		roundRows = append(roundRows, row)
	}

	totalSection := scoreSection([]string{"Total"}, [][]int{totals})
	byCategory := renderScoreTable(header, scoreSection(catLabels, catRows), totalSection)
	header[0] = "Per round"
	byRound := renderScoreTable(header, scoreSection(roundLabels, roundRows), totalSection)
	// Top-align the tables; hjoin centres them vertically.
	for len(byRound) < len(byCategory) {
		byRound = append(byRound, "")
	}
	for len(byCategory) < len(byRound) {
		byCategory = append(byCategory, "")
	}

	screen := vjoin(1,
		block(strings.Split(banner, "\n")),
		block{rule("Final results", 40)},
		hjoin(tableGap, byCategory, byRound),
		renderRanking(totals, result.Winner, ui.self),
	)
	ui.printLines(screen)
	fmt.Fprintln(ui.out)
}

// renderRanking lists players from highest to lowest score, ties sharing a
// place, with the tie-broken winner first among equals.
func renderRanking(totals []int, winner, self int) block {
	order := make([]int, len(totals))
	for i := range order {
		order[i] = i
	}
	slices.SortStableFunc(order, func(a, b int) int {
		if c := cmp.Compare(totals[b], totals[a]); c != 0 {
			return c
		}
		return cmp.Compare(boolInt(b == winner), boolInt(a == winner))
	})

	lines := block{rule("Ranking", 40)}
	place := 0
	for i, p := range order {
		if i == 0 || totals[p] != totals[order[i-1]] {
			place = i + 1
		}
		line := fmt.Sprintf("%d. %-10s %4d points", place, playerName(p, self), totals[p])
		if p == winner {
			line = ansiBold + ansiGreen + line + "  ★ winner" + ansiReset
		}
		lines = append(lines, line)
	}

	verdict := fmt.Sprintf("Player %d wins the game!", winner+1)
	if winner == self {
		verdict = "You win the game!"
	}
	return append(lines, "", ansiBold+verdict+ansiReset)
}

// scoreRow is a labelled row of a score table.
type scoreRow struct {
	label  string
	values []int
}

func scoreSection(labels []string, rows [][]int) []scoreRow {
	out := make([]scoreRow, len(rows))
	for i := range rows {
		out[i] = scoreRow{labels[i], rows[i]}
	}
	return out
}

// renderScoreTable draws a boxed table with one column per player and a
// separator line between sections. The highest positive score in each row is
// highlighted.
func renderScoreTable(header []string, sections ...[]scoreRow) block {
	widths := make([]int, len(header))
	for i, h := range header {
		widths[i] = visibleWidth(h)
	}
	for _, sec := range sections {
		for _, r := range sec {
			widths[0] = max(widths[0], visibleWidth(r.label))
			for i, v := range r.values {
				widths[i+1] = max(widths[i+1], len(fmt.Sprint(v)))
			}
		}
	}

	line := func(l, m, r string) string {
		parts := make([]string, len(widths))
		for i, w := range widths {
			parts[i] = strings.Repeat("─", w+2)
		}
		return l + strings.Join(parts, m) + r
	}
	row := func(cells []string) string {
		parts := make([]string, len(cells))
		for i, c := range cells {
			pad := strings.Repeat(" ", widths[i]-visibleWidth(c))
			if i == 0 {
				parts[i] = " " + c + pad + " "
			} else {
				parts[i] = " " + pad + c + " "
			}
		}
		return "│" + strings.Join(parts, "│") + "│"
	}

	head := make([]string, len(header))
	for i, h := range header {
		head[i] = ansiBold + h + ansiReset
	}
	out := block{line("╭", "┬", "╮"), row(head)}
	for _, sec := range sections {
		if len(sec) == 0 {
			continue
		}
		out = append(out, line("├", "┼", "┤"))
		for _, r := range sec {
			best := 0
			for _, v := range r.values {
				best = max(best, v)
			}
			cells := []string{r.label}
			for _, v := range r.values {
				s := fmt.Sprint(v)
				if v == best && best > 0 {
					s = ansiGreen + s + ansiReset
				}
				cells = append(cells, s)
			}
			out = append(out, row(cells))
		}
	}
	return append(out, line("╰", "┴", "╯"))
}

// presentCategories returns the categories any player has a score in, in
// scoreOrder, followed by any unknown ones.
func presentCategories(scores map[int]map[core.Category]int) []core.Category {
	seen := make(map[core.Category]bool)
	for _, cats := range scores {
		for cat := range cats {
			seen[cat] = true
		}
	}
	var out []core.Category
	for _, cat := range scoreOrder {
		if seen[cat] {
			out = append(out, cat)
			delete(seen, cat)
		}
	}
	rest := make([]core.Category, 0, len(seen))
	for cat := range seen {
		rest = append(rest, cat)
	}
	slices.Sort(rest)
	return append(out, rest...)
}

func sortedKeys[V any](m map[int]V) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Package shared holds scoring building blocks used by several editions.
package shared

import (
	"slices"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

// ScoringFunc adapts a plain function to core.ScoringRule.
type ScoringFunc func(p int, all []*core.PlayerState) int

func (s ScoringFunc) Score(p int, a []*core.PlayerState) int { return s(p, a) }

// ScoreThrowDifference scores the difference between the Throw (first) and
// Catch (last) card values.
func ScoreThrowDifference(p int, all []*core.PlayerState) int {
	draft := all[p].Draft
	if len(draft) == 0 {
		return 0
	}
	diff := draft[0].Value - draft[len(draft)-1].Value
	if diff < 0 {
		return -diff
	}
	return diff
}

// AllVisited returns the sites a player visited in earlier rounds together
// with the sites in the current draft.
func AllVisited(state *core.PlayerState) []byte {
	visited := slices.Clone(state.VisitedSites)
	for _, card := range state.Draft {
		if !slices.Contains(visited, card.Site) {
			visited = append(visited, card.Site)
		}
	}
	return visited
}

// ScoreSites gives one point per drafted site not visited in an earlier round.
func ScoreSites(p int, all []*core.PlayerState) int {
	var score int
	var seen []byte
	for _, card := range all[p].Draft {
		if !slices.Contains(all[p].VisitedSites, card.Site) && !slices.Contains(seen, card.Site) {
			seen = append(seen, card.Site)
			score++
		}
	}
	return score
}

// ScoreRegions awards bonus[i] for clusters[i] when the player has now visited
// every site in it and nobody had completed it in an earlier round. Players
// completing the same region in the same round all get the bonus.
func ScoreRegions(p int, all []*core.PlayerState, clusters [][]byte, bonus []int) int {
	visited := AllVisited(all[p])
	var score int
REGION:
	for i, set := range clusters {
		for _, state := range all {
			if core.SliceIsSubset(set, state.VisitedSites) {
				continue REGION
			}
		}
		if core.SliceIsSubset(set, visited) {
			score += bonus[i]
		}
	}
	return score
}

// ScorePairs scores each pair of identical symbols in a category with the
// symbol's value.
func ScorePairs(p int, all []*core.PlayerState, cat core.Category) int {
	var score int
	for sym, num := range core.FilterSymbolCategoryWithCount(all[p].Draft, cat) {
		score += sym.Value * (num / 2)
	}
	return score
}

// SymbolSum adds up the values of all symbols of a category in the draft.
func SymbolSum(draft []*core.Card, cat core.Category) int {
	var sum int
	for _, card := range draft {
		if sym, ok := card.Symbols[cat]; ok {
			sum += sym.Value
		}
	}
	return sum
}

// ScoreActivityTable scores the chosen activity by how many of it the draft
// holds, looked up in table. Each activity scores only once per game.
func ScoreActivityTable(p int, all []*core.PlayerState, table []int) int {
	choice := all[p].ChosenActivity
	if choice == nil || slices.Contains(all[p].ScoredActivities, choice) {
		return 0
	}
	var count int
	for _, card := range all[p].Draft {
		if card.Symbols[core.Activity] == choice {
			count++
		}
	}
	return table[min(count, len(table)-1)]
}

// RankedBonus scores a once-per-game goal where the first round anyone
// reaches it pays tiers[0], the next round tiers[1] and so on. Players reaching
// it in the same round get the same tier. When the tiers run out, later
// players get the last tier if repeatLast is set, nothing otherwise.
//
// Which tiers are taken is read from earlier rounds' entries of cat on the
// score sheets, so cat must hold nothing but this bonus.
func RankedBonus(p int, all []*core.PlayerState, cat core.Category, reached func(visited []byte) bool, tiers []int, repeatLast bool) int {
	if reached(all[p].VisitedSites) || !reached(AllVisited(all[p])) {
		return 0
	}
	// The player's own entry for this round is not written yet, so its length
	// is the number of rounds already scored.
	rounds := len(all[p].ScoreSheet.CategoryScores[cat])
	taken := 0
	for r := range rounds {
		for _, state := range all {
			if scores := state.ScoreSheet.CategoryScores[cat]; r < len(scores) && scores[r] > 0 {
				taken++
				break
			}
		}
	}
	switch {
	case taken < len(tiers):
		return tiers[taken]
	case repeatLast && len(tiers) > 0:
		return tiers[len(tiers)-1]
	default:
		return 0
	}
}

// GetWinner returns the seat with the highest total score. Ties are decided by
// the highest total in each of cats in turn; if still tied, the lowest seat
// wins. It returns -1 for no score sheets.
func GetWinner(sheets []*core.ScoreSheet, cats ...core.Category) int {
	best := -1
	for id, sheet := range sheets {
		if best < 0 || better(sheet, sheets[best], cats) {
			best = id
		}
	}
	return best
}

// better reports whether a beats b on total score, then on cats in turn.
func better(a, b *core.ScoreSheet, cats []core.Category) bool {
	if ta, tb := a.Total(), b.Total(); ta != tb {
		return ta > tb
	}
	for _, cat := range cats {
		ta, tb := a.CategoryTotal(cat), b.CategoryTotal(cat)
		if ta != tb {
			return ta > tb
		}
	}
	return false
}

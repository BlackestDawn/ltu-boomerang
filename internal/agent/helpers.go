package agent

import (
	"slices"
	"strings"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

// activityOptions lists the activities in hand that have not been scored yet,
// sorted by name, together with how many of each the hand holds.
func activityOptions(hand []*core.Card, alreadyScored []*core.Symbol) ([]*core.Symbol, []int) {
	var options []*core.Symbol
	all := core.FilterSymbolCategoryWithCount(hand, core.Activity)
	for symbol := range all {
		if !slices.Contains(alreadyScored, symbol) {
			options = append(options, symbol)
		}
	}
	slices.SortFunc(options, func(a, b *core.Symbol) int { return strings.Compare(a.Name, b.Name) })

	counts := make([]int, len(options))
	for i, o := range options {
		counts[i] = all[o]
	}
	return options, counts
}

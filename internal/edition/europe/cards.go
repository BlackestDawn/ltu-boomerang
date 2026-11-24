package europe

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

// The base game's 28 cards. The two promo cards (Bulgaria and Romania) are
// left out as there are no published rules for them.
func buildCards() []*core.Card {
	return []*core.Card{
		{Name: "Iceland", Site: 'A', Region: regionNames[0], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[planes],
			core.Activity: cardSymbols[music],
		}},
		{Name: "Norway", Site: 'B', Region: regionNames[0], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[beer],
			core.Activity:   cardSymbols[architecture],
		}},
		{Name: "Sweden", Site: 'C', Region: regionNames[0], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[boats],
			core.Activity: cardSymbols[nature],
		}},
		{Name: "Finland", Site: 'D', Region: regionNames[0], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[cars],
			core.Activity: cardSymbols[nature],
		}},
		{Name: "Denmark", Site: 'E', Region: regionNames[0], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[cars],
			core.Activity: cardSymbols[art],
		}},
		{Name: "Estonia", Site: 'F', Region: regionNames[1], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[boats],
			core.Activity: cardSymbols[music],
		}},
		{Name: "Latvia", Site: 'G', Region: regionNames[1], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[cars],
			core.Activity: cardSymbols[music],
		}},
		{Name: "Lithuania", Site: 'H', Region: regionNames[1], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[trains],
			core.Activity: cardSymbols[art],
		}},
		{Name: "Ireland", Site: 'I', Region: regionNames[2], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[beer],
			core.Activity:   cardSymbols[nature],
		}},
		{Name: "United Kingdom", Site: 'J', Region: regionNames[2], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[cars],
			core.Activity: cardSymbols[architecture],
		}},
		{Name: "Belgium", Site: 'K', Region: regionNames[3], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[beer],
			core.Animal:     cardSymbols[trains],
		}},
		{Name: "Netherlands", Site: 'L', Region: regionNames[3], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[cheese],
			core.Animal:     cardSymbols[bicycles],
		}},
		{Name: "Luxembourg", Site: 'M', Region: regionNames[3], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[cheese],
			core.Activity:   cardSymbols[art],
		}},
		{Name: "Germany", Site: 'N', Region: regionNames[4], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[beer],
			core.Animal:     cardSymbols[trains],
		}},
		{Name: "Switzerland", Site: 'O', Region: regionNames[4], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[cheese],
			core.Activity:   cardSymbols[nature],
		}},
		{Name: "Austria", Site: 'P', Region: regionNames[4], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[trains],
			core.Activity: cardSymbols[music],
		}},
		{Name: "Poland", Site: 'Q', Region: regionNames[5], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[trains],
			core.Activity: cardSymbols[nature],
		}},
		{Name: "Czech Republic", Site: 'R', Region: regionNames[5], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[beer],
			core.Activity:   cardSymbols[architecture],
		}},
		{Name: "Slovakia", Site: 'S', Region: regionNames[5], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[cars],
			core.Activity: cardSymbols[art],
		}},
		{Name: "Hungary", Site: 'T', Region: regionNames[5], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[wine],
			core.Activity:   cardSymbols[art],
		}},
		{Name: "Portugal", Site: 'U', Region: regionNames[6], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[spirits],
			core.Animal:     cardSymbols[bicycles],
		}},
		{Name: "Spain", Site: 'V', Region: regionNames[6], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[cheese],
			core.Activity:   cardSymbols[music],
		}},
		{Name: "France", Site: 'W', Region: regionNames[6], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[wine],
			core.Activity:   cardSymbols[art],
		}},
		{Name: "Italy", Site: 'X', Region: regionNames[6], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[wine],
			core.Activity:   cardSymbols[architecture],
		}},
		{Name: "Slovenia", Site: 'Y', Region: regionNames[6], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[bicycles],
			core.Activity: cardSymbols[nature],
		}},
		{Name: "Croatia", Site: 'Z', Region: regionNames[6], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[boats],
			core.Activity: cardSymbols[music],
		}},
		{Name: "Malta", Site: '@', Region: regionNames[6], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[planes],
			core.Activity: cardSymbols[architecture],
		}},
		{Name: "Greece", Site: '#', Region: regionNames[6], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[bicycles],
			core.Activity: cardSymbols[architecture],
		}},
	}
}

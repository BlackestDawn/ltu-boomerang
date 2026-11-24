package usa

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

// The base game's 28 cards. The two promo cards (Honolulu and Anchorage) are
// left out as there are no published rules for them.
func buildCards() []*core.Card {
	return []*core.Card{
		{Name: "Seattle", Site: 'A', Region: regionNames[0], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[trout],
			core.Activity: cardSymbols[dining],
		}},
		{Name: "Portland, Oregon", Site: 'B', Region: regionNames[0], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[caps],
			core.Animal:     cardSymbols[eagles],
		}},
		{Name: "San Francisco", Site: 'C', Region: regionNames[0], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[trout],
			core.Activity: cardSymbols[dining],
		}},
		{Name: "Los Angeles", Site: 'D', Region: regionNames[0], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[mailboxes],
			core.Activity:   cardSymbols[sightseeing],
		}},
		{Name: "Great Falls", Site: 'E', Region: regionNames[1], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[trout],
			core.Activity: cardSymbols[hiking],
		}},
		{Name: "Las Vegas", Site: 'F', Region: regionNames[1], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[jerseys],
			core.Activity:   cardSymbols[dining],
		}},
		{Name: "Salt Lake City", Site: 'G', Region: regionNames[1], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[mailboxes],
			core.Animal:     cardSymbols[eagles],
		}},
		{Name: "Denver", Site: 'H', Region: regionNames[1], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[bears],
			core.Activity: cardSymbols[sports],
		}},
		{Name: "Phoenix", Site: 'I', Region: regionNames[2], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[mailboxes],
			core.Activity:   cardSymbols[dining],
		}},
		{Name: "Albuquerque", Site: 'J', Region: regionNames[2], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[cattle],
			core.Activity: cardSymbols[hiking],
		}},
		{Name: "Dallas", Site: 'K', Region: regionNames[2], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[cattle],
			core.Activity: cardSymbols[sports],
		}},
		{Name: "Tulsa", Site: 'L', Region: regionNames[2], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[caps],
			core.Animal:     cardSymbols[trout],
		}},
		{Name: "Rapid City", Site: 'M', Region: regionNames[3], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[bears],
			core.Activity: cardSymbols[hiking],
		}},
		{Name: "Minneapolis", Site: 'N', Region: regionNames[3], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[bears],
			core.Activity: cardSymbols[hiking],
		}},
		{Name: "Chicago", Site: 'O', Region: regionNames[3], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[jerseys],
			core.Activity:   cardSymbols[sightseeing],
		}},
		{Name: "St. Louis", Site: 'P', Region: regionNames[3], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[cattle],
			core.Activity: cardSymbols[sports],
		}},
		{Name: "Nashville", Site: 'Q', Region: regionNames[4], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[cattle],
			core.Activity: cardSymbols[sightseeing],
		}},
		{Name: "New Orleans", Site: 'R', Region: regionNames[4], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[mailboxes],
			core.Activity:   cardSymbols[dining],
		}},
		{Name: "Atlanta", Site: 'S', Region: regionNames[4], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[caps],
			core.Activity:   cardSymbols[sports],
		}},
		{Name: "Miami", Site: 'T', Region: regionNames[4], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[trout],
			core.Activity: cardSymbols[sightseeing],
		}},
		{Name: "Columbus", Site: 'U', Region: regionNames[5], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[caps],
			core.Animal:     cardSymbols[cattle],
		}},
		{Name: "Pittsburgh", Site: 'V', Region: regionNames[5], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[bears],
			core.Activity: cardSymbols[sports],
		}},
		{Name: "Washington D.C.", Site: 'W', Region: regionNames[5], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[flags],
			core.Activity:   cardSymbols[sightseeing],
		}},
		{Name: "Atlantic City", Site: 'X', Region: regionNames[5], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[mailboxes],
			core.Activity:   cardSymbols[sightseeing],
		}},
		{Name: "Burlington", Site: 'Y', Region: regionNames[6], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[eagles],
			core.Activity: cardSymbols[hiking],
		}},
		{Name: "New York City", Site: 'Z', Region: regionNames[6], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[jerseys],
			core.Activity:   cardSymbols[dining],
		}},
		{Name: "Portland, Maine", Site: '@', Region: regionNames[6], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[bigfoot],
			core.Activity: cardSymbols[hiking],
		}},
		{Name: "Boston", Site: '#', Region: regionNames[6], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[trout],
			core.Activity: cardSymbols[sports],
		}},
	}
}

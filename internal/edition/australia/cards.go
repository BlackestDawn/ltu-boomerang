package australia

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

// {Name: "",Site: "",Region: regions[""], Symbols: map[core.Category]*core.Symbol{}},
func buildCards() []*core.Card {
	return []*core.Card{
		{Name: "The Bungle Bungles", Site: 'A', Region: regionNames[0], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[leaves],
			core.Activity:   cardSymbols[indigenous],
		}},
		{Name: "The Pinnacles", Site: 'B', Region: regionNames[0], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[kangaroos],
			core.Activity: cardSymbols[sightseeing],
		}},
		{Name: "Margaret River", Site: 'C', Region: regionNames[0], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[shells],
			core.Animal:     cardSymbols[kangaroos],
		}},
		{Name: "Kalbarri National Park", Site: 'D', Region: regionNames[0], Value: 1, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[wildflowers],
			core.Activity:   cardSymbols[bushwalking],
		}},
		{Name: "Uluru", Site: 'E', Region: regionNames[1], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[emus],
			core.Activity: cardSymbols[indigenous],
		}},
		{Name: "Kakadu National Park", Site: 'F', Region: regionNames[1], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[wombats],
			core.Activity: cardSymbols[sightseeing],
		}},
		{Name: "Nitmiluk National Park", Site: 'G', Region: regionNames[1], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[shells],
			core.Animal:     cardSymbols[platypuses],
		}},
		{Name: "King's Canyon", Site: 'H', Region: regionNames[1], Value: 4, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[koalas],
			core.Activity: cardSymbols[swimming],
		}},
		{Name: "The Great Barrier Reef", Site: 'I', Region: regionNames[2], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[wildflowers],
			core.Activity:   cardSymbols[sightseeing],
		}},
		{Name: "The Whitsundays", Site: 'J', Region: regionNames[2], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[kangaroos],
			core.Activity: cardSymbols[indigenous],
		}},
		{Name: "Daintree Rainforest", Site: 'K', Region: regionNames[2], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[souvenirs],
			core.Activity:   cardSymbols[bushwalking],
		}},
		{Name: "Surfers Paradise", Site: 'L', Region: regionNames[2], Value: 6, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[wildflowers],
			core.Activity:   cardSymbols[swimming],
		}},
		{Name: "Barossa Valley", Site: 'M', Region: regionNames[3], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[koalas],
			core.Activity: cardSymbols[bushwalking],
		}},
		{Name: "Lake Eyre", Site: 'N', Region: regionNames[3], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[emus],
			core.Activity: cardSymbols[swimming],
		}},
		{Name: "Kangaroo Island", Site: 'O', Region: regionNames[3], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[kangaroos],
			core.Activity: cardSymbols[bushwalking],
		}},
		{Name: "Mount Gambier", Site: 'P', Region: regionNames[3], Value: 3, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[wildflowers],
			core.Activity:   cardSymbols[sightseeing],
		}},
		{Name: "Blue Mountains", Site: 'Q', Region: regionNames[4], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[wombats],
			core.Activity: cardSymbols[indigenous],
		}},
		{Name: "Sydney Harbour", Site: 'R', Region: regionNames[4], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[emus],
			core.Activity: cardSymbols[sightseeing],
		}},
		{Name: "Bondi Beach", Site: 'S', Region: regionNames[4], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[wombats],
			core.Activity: cardSymbols[swimming],
		}},
		{Name: "Hunter Valley", Site: 'T', Region: regionNames[4], Value: 5, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[emus],
			core.Activity: cardSymbols[bushwalking],
		}},
		{Name: "Melbourne", Site: 'U', Region: regionNames[5], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[wombats],
			core.Activity: cardSymbols[bushwalking],
		}},
		{Name: "The MCG", Site: 'V', Region: regionNames[5], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[leaves],
			core.Activity:   cardSymbols[indigenous],
		}},
		{Name: "Twelve Apostles", Site: 'W', Region: regionNames[5], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[shells],
			core.Activity:   cardSymbols[swimming],
		}},
		{Name: "Royal Exhibition Building", Site: 'X', Region: regionNames[5], Value: 2, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[leaves],
			core.Animal:     cardSymbols[platypuses],
		}},
		{Name: "Salamanca Markets", Site: 'Y', Region: regionNames[6], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[leaves],
			core.Animal:     cardSymbols[emus],
		}},
		{Name: "Mount Wellington", Site: 'Z', Region: regionNames[6], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[koalas],
			core.Activity: cardSymbols[sightseeing],
		}},
		{Name: "Port Arthur", Site: '*', Region: regionNames[6], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Collection: cardSymbols[leaves],
			core.Activity:   cardSymbols[indigenous],
		}},
		{Name: "Richmond", Site: '-', Region: regionNames[6], Value: 7, Symbols: map[core.Category]*core.Symbol{
			core.Animal:   cardSymbols[kangaroos],
			core.Activity: cardSymbols[swimming],
		}},
	}
}

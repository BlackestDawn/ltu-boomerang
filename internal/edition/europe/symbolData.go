package europe

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

type symbolID int

const (
	beer symbolID = iota
	cheese
	wine
	spirits
	cars
	trains
	bicycles
	boats
	planes
	art
	architecture
	music
	nature
)

var cardSymbols = map[symbolID]*core.Symbol{
	beer:         {Name: "Beer", Value: 1},
	cheese:       {Name: "Cheese", Value: 2},
	wine:         {Name: "Wine", Value: 3},
	spirits:      {Name: "Spirits", Value: 5},
	cars:         {Name: "Cars", Value: 3},
	trains:       {Name: "Trains", Value: 4},
	bicycles:     {Name: "Bicycles", Value: 5},
	boats:        {Name: "Boats", Value: 7},
	planes:       {Name: "Planes", Value: 9},
	art:          {Name: "Art", Value: 0},
	architecture: {Name: "Architecture", Value: 0},
	music:        {Name: "Music", Value: 0},
	nature:       {Name: "Natural Wonders", Value: 0},
}

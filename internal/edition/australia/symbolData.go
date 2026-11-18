package australia

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

type symbolID int

const (
	leaves symbolID = iota
	wildflowers
	shells
	souvenirs
	kangaroos
	emus
	wombats
	koalas
	platypuses
	swimming
	bushwalking
	indigenous
	sightseeing
)

var cardSymbols = map[symbolID]*core.Symbol{
	leaves:      {Name: "Leaves", Value: 1},
	wildflowers: {Name: "Wildflowers", Value: 2},
	shells:      {Name: "Shells", Value: 3},
	souvenirs:   {Name: "Souvenirs", Value: 5},
	kangaroos:   {Name: "Kangaroos", Value: 3},
	emus:        {Name: "Emus", Value: 4},
	wombats:     {Name: "Wombats", Value: 5},
	koalas:      {Name: "Koalas", Value: 7},
	platypuses:  {Name: "Platypuses", Value: 9},
	swimming:    {Name: "Swimming", Value: 0},
	bushwalking: {Name: "Bushwalking", Value: 0},
	indigenous:  {Name: "Indigenous Culture", Value: 0},
	sightseeing: {Name: "Sightseeing", Value: 0},
}

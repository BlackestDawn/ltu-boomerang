package usa

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

type symbolID int

const (
	mailboxes symbolID = iota
	caps
	jerseys
	flags
	trout
	cattle
	bears
	eagles
	bigfoot
	sightseeing
	sports
	hiking
	dining
)

var cardSymbols = map[symbolID]*core.Symbol{
	mailboxes:   {Name: "Mail Boxes", Value: 1},
	caps:        {Name: "Baseball Caps", Value: 2},
	jerseys:     {Name: "Football Jerseys", Value: 3},
	flags:       {Name: "US Flags", Value: 5},
	trout:       {Name: "Trout", Value: 3},
	cattle:      {Name: "Cattle", Value: 4},
	bears:       {Name: "Grizzly Bears", Value: 5},
	eagles:      {Name: "Bald Eagles", Value: 7},
	bigfoot:     {Name: "Big Foot", Value: 9},
	sightseeing: {Name: "Sightseeing", Value: 0},
	sports:      {Name: "Sporting Events", Value: 0},
	hiking:      {Name: "Hiking", Value: 0},
	dining:      {Name: "Dining", Value: 0},
}

package australia

import (
	"testing"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

// expectedCards is the card table from the exam (page 2), with "New South
// Wales" spelled correctly. Empty strings mean the card has no symbol in that
// category.
var expectedCards = []struct {
	name                         string
	site                         byte
	region                       string
	value                        int
	collection, animal, activity string
}{
	{"The Bungle Bungles", 'A', "Western Australia", 1, "Leaves", "", "Indigenous Culture"},
	{"The Pinnacles", 'B', "Western Australia", 1, "", "Kangaroos", "Sightseeing"},
	{"Margaret River", 'C', "Western Australia", 1, "Shells", "Kangaroos", ""},
	{"Kalbarri National Park", 'D', "Western Australia", 1, "Wildflowers", "", "Bushwalking"},
	{"Uluru", 'E', "Northern Territory", 4, "", "Emus", "Indigenous Culture"},
	{"Kakadu National Park", 'F', "Northern Territory", 4, "", "Wombats", "Sightseeing"},
	{"Nitmiluk National Park", 'G', "Northern Territory", 4, "Shells", "Platypuses", ""},
	{"King's Canyon", 'H', "Northern Territory", 4, "", "Koalas", "Swimming"},
	{"The Great Barrier Reef", 'I', "Queensland", 6, "Wildflowers", "", "Sightseeing"},
	{"The Whitsundays", 'J', "Queensland", 6, "", "Kangaroos", "Indigenous Culture"},
	{"Daintree Rainforest", 'K', "Queensland", 6, "Souvenirs", "", "Bushwalking"},
	{"Surfers Paradise", 'L', "Queensland", 6, "Wildflowers", "", "Swimming"},
	{"Barossa Valley", 'M', "South Australia", 3, "", "Koalas", "Bushwalking"},
	{"Lake Eyre", 'N', "South Australia", 3, "", "Emus", "Swimming"},
	{"Kangaroo Island", 'O', "South Australia", 3, "", "Kangaroos", "Bushwalking"},
	{"Mount Gambier", 'P', "South Australia", 3, "Wildflowers", "", "Sightseeing"},
	{"Blue Mountains", 'Q', "New South Wales", 5, "", "Wombats", "Indigenous Culture"},
	{"Sydney Harbour", 'R', "New South Wales", 5, "", "Emus", "Sightseeing"},
	{"Bondi Beach", 'S', "New South Wales", 5, "", "Wombats", "Swimming"},
	{"Hunter Valley", 'T', "New South Wales", 5, "", "Emus", "Bushwalking"},
	{"Melbourne", 'U', "Victoria", 2, "", "Wombats", "Bushwalking"},
	{"The MCG", 'V', "Victoria", 2, "Leaves", "", "Indigenous Culture"},
	{"Twelve Apostles", 'W', "Victoria", 2, "Shells", "", "Swimming"},
	{"Royal Exhibition Building", 'X', "Victoria", 2, "Leaves", "Platypuses", ""},
	{"Salamanca Markets", 'Y', "Tasmania", 7, "Leaves", "Emus", ""},
	{"Mount Wellington", 'Z', "Tasmania", 7, "", "Koalas", "Sightseeing"},
	{"Port Arthur", '*', "Tasmania", 7, "Leaves", "", "Indigenous Culture"},
	{"Richmond", '-', "Tasmania", 7, "", "Kangaroos", "Swimming"},
}

// expectedSymbolValues are the point values from requirement 10c and 10d.
var expectedSymbolValues = map[string]int{
	"Leaves": 1, "Wildflowers": 2, "Shells": 3, "Souvenirs": 5,
	"Kangaroos": 3, "Emus": 4, "Wombats": 5, "Koalas": 7, "Platypuses": 9,
}

func symbolName(c *core.Card, cat core.Category) string {
	if s, ok := c.Symbols[cat]; ok {
		return s.Name
	}
	return ""
}

// Requirement 2: the deck consists of 28 cards.
func TestDeckHas28Cards(t *testing.T) {
	if got := len(New().Cards()); got != 28 {
		t.Fatalf("deck has %d cards, want 28", got)
	}
}

// Requirement 2a-c: every card matches the exam's card table.
func TestCardsMatchTable(t *testing.T) {
	bySite := map[byte]*core.Card{}
	for _, c := range New().Cards() {
		if _, dup := bySite[c.Site]; dup {
			t.Errorf("site %c appears twice", c.Site)
		}
		bySite[c.Site] = c
	}

	for _, want := range expectedCards {
		c, ok := bySite[want.site]
		if !ok {
			t.Errorf("no card for site %c (%s)", want.site, want.name)
			continue
		}
		if c.Name != want.name || c.Region != want.region || c.Value != want.value {
			t.Errorf("site %c: got %q/%q/%d, want %q/%q/%d",
				want.site, c.Name, c.Region, c.Value, want.name, want.region, want.value)
		}
		if got := symbolName(c, core.Collection); got != want.collection {
			t.Errorf("%s collection: got %q, want %q", want.name, got, want.collection)
		}
		if got := symbolName(c, core.Animal); got != want.animal {
			t.Errorf("%s animal: got %q, want %q", want.name, got, want.animal)
		}
		if got := symbolName(c, core.Activity); got != want.activity {
			t.Errorf("%s activity: got %q, want %q", want.name, got, want.activity)
		}
	}
}

// Requirement 2c: each card has exactly two symbols.
func TestEveryCardHasTwoSymbols(t *testing.T) {
	for _, c := range New().Cards() {
		if len(c.Symbols) != 2 {
			t.Errorf("%s has %d symbols, want 2", c.Name, len(c.Symbols))
		}
	}
}

// Requirement 10c/10d: symbols carry the point values from the rules.
func TestSymbolValues(t *testing.T) {
	for _, c := range New().Cards() {
		for _, cat := range []core.Category{core.Collection, core.Animal} {
			s, ok := c.Symbols[cat]
			if !ok {
				continue
			}
			if want := expectedSymbolValues[s.Name]; s.Value != want {
				t.Errorf("%s on %s worth %d, want %d", s.Name, c.Name, s.Value, want)
			}
		}
	}
}

// Requirement 10b: the region clusters used for bonuses match the regions
// printed on the cards.
func TestRegionClustersMatchCards(t *testing.T) {
	regionOf := map[byte]string{}
	for _, c := range New().Cards() {
		regionOf[c.Site] = c.Region
	}
	if len(regionClusters) != 7 {
		t.Fatalf("%d regions, want 7", len(regionClusters))
	}
	for i, cluster := range regionClusters {
		if len(cluster) != 4 {
			t.Errorf("region %s has %d sites, want 4", regionNames[i], len(cluster))
		}
		for _, site := range cluster {
			if regionOf[site] != regionNames[i] {
				t.Errorf("site %c is in cluster %s but card says %s", site, regionNames[i], regionOf[site])
			}
		}
	}
}

// Shuffling the slice returned by Cards must not reorder the edition's own
// card list.
func TestCardsReturnsCopy(t *testing.T) {
	ed := New()
	first := ed.Cards()
	first[0], first[1] = first[1], first[0]
	if ed.Cards()[0].Site != 'A' {
		t.Error("modifying the returned slice changed the edition's cards")
	}
}

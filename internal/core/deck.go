package core

import (
	"math/rand"
	"slices"
)

type Deck struct {
	Cards        []*Card
	lastDealtPos int
}

func (d *Deck) Shuffle(rng *rand.Rand) {
	d.lastDealtPos = 0
	rng.Shuffle(len(d.Cards), func(i, j int) { d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i] })
}

func (d *Deck) Deal(n int) []*Card {
	cards := slices.Clone(d.Cards[d.lastDealtPos : d.lastDealtPos+n])
	d.lastDealtPos += n
	return cards
}

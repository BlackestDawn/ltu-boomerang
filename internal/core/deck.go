package core

import "math/rand"

type Deck struct {
	Cards []Card
}

func (d *Deck) Shuffler(rng *rand.Rand) {}

func (d *Deck) Deal(n int) []Card {
	return nil
}

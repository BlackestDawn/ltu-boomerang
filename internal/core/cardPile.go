package core

import "slices"

type CardPile struct {
	Cards []*Card
}

func (c *CardPile) Add(card *Card) {
	c.Cards = append(c.Cards, card)
}

func (c *CardPile) Set(cards []*Card) {
	c.Cards = cards
}

func (c *CardPile) Remove(card *Card) {
	c.Cards = slices.DeleteFunc(c.Cards, func(cc *Card) bool { return cc == card })
}

func (c *CardPile) Empty() {
	c.Cards = []*Card{}
}

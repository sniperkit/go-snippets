package main

import (
	"fmt"
	"github.com/aita/engi"
	"github.com/gophergala/gunosy/daihinmin"
)

type CardSprite struct {
	*engi.Sprite
	daihinmin.Card
	selected bool
}

func NewCardSprite(card daihinmin.Card, x, y float32) *CardSprite {
	texture := engi.Files.Image(cardTextureName(card))
	region := engi.NewRegion(texture, 0, 0, int(texture.Width()), int(texture.Height()))
	sprite := &CardSprite{
		engi.NewSprite(region, x, y),
		card,
		false,
	}
	return sprite
}

func (card *CardSprite) HitTest(x, y float32) bool {
	// NOTE: Positionしか考慮しない
	if card.Position.X < x && x < card.Position.X+card.Region.Width() {
		if card.Position.Y < y && y < card.Position.Y+card.Region.Height() {
			return true
		}
	}
	return false
}

func cardTextureName(card daihinmin.Card) string {
	return fmt.Sprintf("%c%d", card.Suit, card.Rank)
}

func (game *Game) loadCardImages() {
	var suits = map[daihinmin.Suit]rune{
		daihinmin.Spade:   's',
		daihinmin.Heart:   'h',
		daihinmin.Diamond: 'd',
		daihinmin.Club:    'c',
	}
	for rank := 1; rank < 14; rank++ {
		for suit, c := range suits {
			card := daihinmin.Card{suit, daihinmin.Rank(rank)}
			engi.Files.Add(
				cardTextureName(card),
				fmt.Sprintf("data/cards/%c%02d.png", c, rank),
			)
		}
	}
}

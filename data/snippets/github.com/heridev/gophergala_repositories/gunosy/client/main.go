package main

import (
	"fmt"
	"github.com/aita/engi"
	"github.com/gophergala/gunosy/daihinmin"
)

type Game struct {
	*engi.Game
	batch *engi.Batch
	cards []*CardSprite
	mouse bool
}

func (game *Game) HitTest(x, y float32) {
	// カードの選択
	for index := range game.cards {
		sprite := game.cards[len(game.cards)-index-1]
		if sprite.HitTest(x, y) {
			sprite.selected = !sprite.selected
			if sprite.selected {
				sprite.Position.Y -= 30
			} else {
				sprite.Position.Y += 30
			}

			fmt.Printf("%c%d: %s\n", sprite.Suit, sprite.Rank, sprite.selected)
			break
		}
	}
}

func (game *Game) Preload() {
	engi.Files.Add("gopher", "data/gopher_s.png")
	engi.Files.Add("back", "data/cards/z02.png")
	engi.Files.Add("font", "data/font.png")
	game.loadCardImages()
	game.batch = engi.NewBatch(engi.Width(), engi.Height())
}

func (game *Game) Setup() {
	engi.SetBg(0x2d3739)
	game.mouse = false

	for n := 1; n < 14; n++ {
		card := daihinmin.Card{daihinmin.Spade, daihinmin.Rank(n)}
		sprite := NewCardSprite(card, float32(140+30*n), 420)
		game.cards = append(game.cards, sprite)
	}
}

func (game *Game) Mouse(x, y float32, action engi.Action) {
	switch action {
	case engi.PRESS:
		game.mouse = true
	case engi.RELEASE:
		if game.mouse {
			game.HitTest(x, y)
		}
		game.mouse = false
	}
}

func (game *Game) Render() {
	game.batch.Begin()
	for _, sprite := range game.cards {
		sprite.Render(game.batch)
	}
	game.batch.End()
}

func main() {
	engi.Open("hello", 800, 600, false, &Game{})
}

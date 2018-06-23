package store

import (
	"fmt"
	"math/rand"

	"github.com/acsellers/words"
)

func (ds DeckScope) AvailableDecks() DeckScope {
	return ds.FullGame().Eq(true).Private().Eq(false)
}

func (d Deck) BuildGame(players int) Game {
	g := Game{
		D:            d,
		CurrentPlays: make(map[string]Card),
	}
	g.GC, _ = d.cached_conn.Card.DeckID().Eq(d.ID).Type().Eq("parent").RetrieveAll()
	// shuffle the cards
	for i := len(g.GC) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		if i != j {
			g.GC[i], g.GC[j] = g.GC[j], g.GC[i]
		}
	}

	var err error
	g.PC, err = d.cached_conn.Card.DeckID().Eq(d.ID).Type().Eq("child").RetrieveAll()
	if err != nil || len(g.PC) == 0 {
		fmt.Println(d.Scope().CardScope().Type().Eq("child").QuerySQL())
	}
	// shuffle the cards
	for i := len(g.PC) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		if i != j {
			g.PC[i], g.PC[j] = g.PC[j], g.PC[i]
		}
	}
	// append some random cards at the end
	for i := 0; i < 50; i++ {
		g.PC = append(g.PC, Card{
			Name: words.QuickDouble(),
			Type: "child",
		})
	}
	g.Hands = make([][]Card, players)
	for i := range g.Hands {
		g.Hands[i] = g.PC[:6]
		g.PC = g.PC[6:]
	}

	g.Players = make(map[string]int)
	return g
}

type Game struct {
	D            Deck
	GC           []Card
	I            int
	PC           []Card
	Hands        [][]Card
	Players      map[string]int
	CurrentPlays map[string]Card
	CurrentJudge int
	RoundCard    Card
	RoundWinner  string
}

func (g Game) CurrentCard() string {
	switch g.D.GameType {
	case "blanks":
		return fmt.Sprintf(g.GC[g.I].Name, "__________")
	case "adjective":
		return g.GC[g.I].Name
	default:
		return "Card not found"
	}
}

func (g Game) CurrentWith(c Card) string {
	switch g.D.GameType {
	case "blanks":
		return fmt.Sprintf(g.GC[g.I].Name, c.Name)
	case "adjective":
		return fmt.Sprintf("%s %s", c.Name, g.GC[g.I].Name)
	default:
		return "Game not found"
	}

}
func (g *Game) AdvanceCard() {
	if g.I == len(g.GC) {
		g.I = 0
	} else {
		g.I++
	}
}

func (g Game) Judging() bool {
	return len(g.CurrentPlays)+1 == len(g.Players)
}

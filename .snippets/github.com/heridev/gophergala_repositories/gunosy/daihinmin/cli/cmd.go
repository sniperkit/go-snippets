package main

import (
	"errors"
	"fmt"
	"github.com/gophergala/gunosy/daihinmin"
	"strconv"
	"strings"
)

func main() {
	game := daihinmin.NewGame()
	game.Join(daihinmin.NewPlayer("Alice"))
	game.Join(daihinmin.NewPlayer("Bob"))
	game.Join(daihinmin.NewPlayer("Charles"))
	game.Join(daihinmin.NewPlayer("Daniel"))
	game.Start()
	printGame(*game)
	fmt.Println()

	for true {
		p := currPlayer(*game)
		fmt.Printf("(%s) ", p.Name)
		for _, c := range p.Hand {
			fmt.Printf("%s", c)
		}
		fmt.Println()
		fmt.Printf("> ")
		var line string
		_, err := fmt.Scanln(&line)
		if err != nil {
			fmt.Printf("Failed to parse: %s\n", err)
			continue
		}
		if line[0] == 'P' {
			events := game.Pass(p)
			fmt.Println(events)
		} else {
			cards, err := readInput(line)
			if err != nil {
				fmt.Printf("Failed to parse [%s]: %s\n", line, err)
				continue
			}
			ok, events := game.Play(p, cards)
			if !ok {
				fmt.Println("Invalid input!")
				fmt.Println(events)
				continue
			}
			fmt.Println(events)
		}
	}
}

func readInput(line string) (daihinmin.Cards, error) {
	cards := make([]daihinmin.Card, 0)
	for _, c := range strings.Split(line, ",") {
		c = strings.TrimSpace(c)
		if len(c) == 0 || len(c) > 4 {
			return nil, errors.New("Invalid input: " + c)
		}
		var s daihinmin.Suit
		var r daihinmin.Rank
		switch c[0] {
		case 'S':
			s = daihinmin.Spade
		case 'H':
			s = daihinmin.Heart
		case 'D':
			s = daihinmin.Diamond
		case 'C':
			s = daihinmin.Club
		case 'J':
			s = daihinmin.Joker
		}

		if s == daihinmin.Joker {
			r = daihinmin.JokerRank
		} else if len(c) < 2 {
			return nil, errors.New("Invalid input: " + c)
		} else {
			switch c[1] {
			case 'A':
				r = daihinmin.Ace
			case 'K':
				r = daihinmin.King
			case 'Q':
				r = daihinmin.Queen
			case 'J':
				r = daihinmin.Jack
			default:
				num, error := strconv.Atoi(c[1:])
				if error != nil {
					return nil, error
				}
				r = daihinmin.Rank(num)
			}
		}

		cards = append(cards, daihinmin.Card{s, r})
	}
	return daihinmin.Cards(cards), nil
}

func currPlayer(g daihinmin.Game) *daihinmin.Player {
	for _, p := range g.Players {
		if p.Number == g.Current {
			return p
		}
	}
	return nil
}

func printGame(g daihinmin.Game) {
	for _, p := range g.Players {
		m := ""
		if p.Number == g.Current {
			m = "*"
		}
		fmt.Printf("%s(%s) [", m, p.Name)
		for _, c := range p.Hand {
			fmt.Printf("%s", c)
		}
		fmt.Printf("]\n")
	}
}

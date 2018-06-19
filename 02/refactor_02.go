package main

import (
	"fmt"
	"strconv"
)

type Player struct {
	Name  string
	Point int
}

type TennisGame struct {
	Player1 *Player
	Player2 *Player
}

func (g TennisGame) AddPoint(player *Player) {
	player.Point += 1
}

func (g TennisGame) GetScore() string {
	var playerName string
	var tempScore string

	if g.Player1.Point < 4 && g.Player2.Point < 4 && !(g.Player1.Point+g.Player2.Point == 6) {
		p := []string{"Love", "Fifteen", "Thirty", "Forty"}
		tempScore = p[g.Player1.Point]
		if g.Player1.Point == g.Player2.Point {
			return tempScore + "-All"
		}
		return tempScore + "-" + p[g.Player2.Point]
	} else {
		if g.Player1.Point == g.Player2.Point {
			return "Deuce"
		}
		if g.Player1.Point > g.Player2.Point {
			playerName = g.Player1.Name
		} else {
			playerName = g.Player2.Name
		}
		if (g.Player1.Point-g.Player2.Point)*(g.Player1.Point-g.Player2.Point) == 1 {
			return "Advantage " + playerName
		}
		return "Win for " + playerName
	}
}

func (g TennisGame) printScore() string {
	text := g.Player1.Name + " " + strconv.Itoa(g.Player1.Point) + " - " + strconv.Itoa(g.Player2.Point) + " " + g.Player2.Name
	return text
}

func main() {
	var player1 = &Player{"Sepry", 0}
	var player2 = &Player{"Haryandi", 0}
	var g = TennisGame{player1, player2}

	g.AddPoint(player1)
	fmt.Println(g.printScore())
	fmt.Println(g.GetScore())
}

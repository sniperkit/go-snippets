package main

import (
	"fmt"
	"strconv"
)

var p = []string{"Love", "Fifteen", "Thirty", "Forty"}
var win string = "Win for "
var adv string = "Advantage "
var deuce string = "Deuce"
var all string = "All"

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
	if g.isWinner() {
		return win + g.leadingPlayer()
	}
	if g.isAdvantage() {
		return adv + g.leadingPlayer()
	}
	if g.isDeuce() {
		return deuce
	}
	return g.generalScore()
}

func (g TennisGame) isDeuce() bool {
	return g.Player1.Point+g.Player2.Point >= 6 && g.isTied()
}

func (g TennisGame) isAdvantage() bool {
	return g.isLeadingByOne() && g.hasPointToWon()
}

func (g TennisGame) isWinner() bool {
	return g.isLeadingByAtLeastTwo() && g.hasPointToWon()
}

func (g TennisGame) isTied() bool {
	return g.Player1.Point == g.Player2.Point
}

func (g TennisGame) leadingPlayer() string {
	playerName := g.Player1.Name
	if g.Player2.Point > g.Player1.Point {
		return g.Player2.Name
	}
	return playerName
}

func (g TennisGame) isLeadingByOne() bool {
	return (g.Player1.Point-g.Player2.Point)*(g.Player1.Point-g.Player2.Point) == 1
}

func (g TennisGame) isLeadingByAtLeastTwo() bool {
	return (g.Player1.Point-g.Player2.Point)*(g.Player1.Point-g.Player2.Point) >= 2
}

func (g TennisGame) hasPointToWon() bool {
	return g.Player1.Point >= 4 || g.Player2.Point >= 4
}

func (g TennisGame) generalScore() string {
	score := p[g.Player1.Point] + "-"
	if g.isTied() {
		return score + all
	} else {
		return score + p[g.Player2.Point]
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
	g.AddPoint(player1)
	g.AddPoint(player1)
	g.AddPoint(player2)
	g.AddPoint(player2)
	fmt.Println(g.printScore())
	fmt.Println(g.GetScore())
}

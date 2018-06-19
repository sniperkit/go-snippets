package main

import (
	"testing"
)

var player1 = &Player{"Sepry", 0}
var player2 = &Player{"Haryandi", 0}
var g = TennisGame{player1, player2}

// reset players point for every test
func Teardown() {
	player1.Point = 0
	player2.Point = 0
}

func TestAllPlayerLove(t *testing.T) {
	t.Log(g.printScore())
	t.Log(g.GetScore())

	Teardown()
}

func TestAllPlayerFifteen(t *testing.T) {
	g.AddPoint(player1)
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	Teardown()
}

func TestFirstPlayerWinFirstGame(t *testing.T) {
	g.AddPoint(player1)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	Teardown()
}

func TestSecondPlayerWinFirstGame(t *testing.T) {
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	Teardown()
}

func TestFirstPlayerAdvantage(t *testing.T) {
	g.AddPoint(player1)
	g.AddPoint(player1)
	g.AddPoint(player1)
	g.AddPoint(player1)

	g.AddPoint(player2)
	g.AddPoint(player2)
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	Teardown()
}

func TestFirstPlayerWinAll(t *testing.T) {
	g.AddPoint(player1)
	g.AddPoint(player1)
	g.AddPoint(player1)
	g.AddPoint(player1)

	g.AddPoint(player2)
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	Teardown()
}

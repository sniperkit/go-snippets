package main

import (
	"fmt"
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

	expected := "Love-All"
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestAllPlayerFifteen(t *testing.T) {
	g.AddPoint(player1)
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	expected := "Fifteen-All"
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestAllPlayerThirty(t *testing.T) {
	g.AddPoint(player1)
	g.AddPoint(player1)
	g.AddPoint(player2)
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	expected := "Thirty-All"
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestAllPlayerDeuce(t *testing.T) {
	g.AddPoint(player1)
	g.AddPoint(player1)
	g.AddPoint(player1)
	g.AddPoint(player2)
	g.AddPoint(player2)
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	expected := "Deuce"
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestFirstPlayerWinFirstGame(t *testing.T) {
	g.AddPoint(player1)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	expected := "Fifteen-Love"
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestFirstPlayerWinSecondGame(t *testing.T) {
	g.AddPoint(player1)
	g.AddPoint(player1)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	expected := "Thirty-Love"
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestFirstPlayerWinThirdGame(t *testing.T) {
	g.AddPoint(player1)
	g.AddPoint(player1)
	g.AddPoint(player1)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	expected := "Forty-Love"
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestSecondPlayerWinFirstGame(t *testing.T) {
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	expected := "Love-Fifteen"
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestSecondPlayerWinSecondGame(t *testing.T) {
	g.AddPoint(player2)
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	expected := "Love-Thirty"
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestSecondPlayerWinThirdGame(t *testing.T) {
	g.AddPoint(player2)
	g.AddPoint(player2)
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	expected := "Love-Forty"
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

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

	expected := "Advantage " + g.Player1.Name
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestSecondPlayerAdvantage(t *testing.T) {
	g.AddPoint(player1)
	g.AddPoint(player1)
	g.AddPoint(player1)

	g.AddPoint(player2)
	g.AddPoint(player2)
	g.AddPoint(player2)
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	expected := "Advantage " + g.Player2.Name
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

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

	expected := "Win for " + g.Player1.Name
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestSecondPlayerWinAll(t *testing.T) {
	g.AddPoint(player1)
	g.AddPoint(player1)

	g.AddPoint(player2)
	g.AddPoint(player2)
	g.AddPoint(player2)
	g.AddPoint(player2)

	t.Log(g.printScore())
	t.Log(g.GetScore())

	expected := "Win for " + g.Player2.Name
	actual := g.GetScore()

	if actual != expected {
		t.Error(fmt.Sprintf("Expected result is %d but instead got %d!", expected, actual))
	}

	Teardown()
}

func TestMain(t *testing.T) {
	main()

	Teardown()
}

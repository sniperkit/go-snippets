package controllers

import (
	"github.com/go_warrior/characters"
	"github.com/go_warrior/game"

	"fmt"
	"time"
)

type LevelController struct {
	Printer *Printer
	Warrior *characters.Warrior
	Board   *game.Board
}

func NewLevel() *LevelController {
	generator := NewBoardGenerator()

	warrior, board, err := generator.Generate("00")
	panicIfError(err)

	return &LevelController{
		Printer: &Printer{
			Board: board,
		},
		Warrior: warrior,
		Board:   board,
	}
}

func (this *LevelController) Start(f UserFunction) {
	defer func() {
		if r := recover(); r != nil {
			endLevel(r)
		}
	}()

	for {
		this.Printer.PrintBoard()

		turn := &TurnController{
			Board: this.Board,
		}

		p, i := turn.getAllCharacters()
		turn.resolveTurns(f, p, i)

		if isLevelSucceded() {
			endLevel("Level successful!")
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func endLevel(cause interface{}) {
	var message string

	switch causeType := cause.(type) {
	case string:
		message = causeType
		break
	case error:
		message = causeType.Error()
		break
	default:
		message = "The game ended for unknown reasons"
	}

	fmt.Println(message)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

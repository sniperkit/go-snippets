package controllers

import (
	"fmt"

	"github.com/go_warrior/game"
)

type Printer struct {
	Board *game.Board
}

func (this *Printer) PrintBoard() {
	fmt.Println("----")

	for i := 0; i < this.Board.Height; i++ {
		fmt.Print("|")

		for j := 0; j < this.Board.Width; j++ {
			key := GenerateKey(j, i)
			currentSpace := this.Board.Spaces[key]

			if currentSpace.Empty() {
				fmt.Print(" ")
			} else {
				fmt.Print(currentSpace.Element.GetSprite())
			}
		}

		fmt.Println("|")
	}

	fmt.Println("----")
}

package main

import (
	"fmt"
)

type CommandHelp struct {
	Line1 string
	Line2 string
	Line3 string
	Line4 string
	Line5 string
}

func (command *CommandHelp) Init() {
	command.Line1 = "adventure\t: Set out on an adventure to find experience, loot, and random battles."
	command.Line2 = "me\t\t: Look at your stats."
	command.Line3 = "log\t\t: View your adventure so far."
	command.Line4 = "quit\t\t: Quit the game :-("
	command.Line5 = "help\t\t: Get command help."

}

func (command *CommandHelp) PrintHelpCommands() {
	fmt.Println()
	fmt.Println("Game Commands")
	fmt.Println("-------------")
	fmt.Println(command.Line1)
	fmt.Println(command.Line2)
	fmt.Println(command.Line3)
	fmt.Println(command.Line4)
	fmt.Println(command.Line5)
	fmt.Println()
}

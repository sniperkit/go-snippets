package main

import (
	"bufio"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"os"
	"os/exec"
	"runtime"
)

func failOnError(err error) {
	if err != nil {
		ct.ChangeColor(ct.Red, true, ct.None, false)
		fmt.Println(err)
		ct.ResetColor()
		os.Exit(1)
	}
}

func main() {
	defer ct.ResetColor()
	ClearScreen()

	ct.ChangeColor(ct.Green, true, ct.None, false)

	fmt.Println()
	fmt.Println("Welcome to shark_sandwich!")
	fmt.Println()

	ConsoleReader := bufio.NewReader(os.Stdin)
	storage, err := NewStorage()
	failOnError(err)

	hero, err := InitGame(ConsoleReader, storage)
	failOnError(err)

	gameWorld := NewGameWorld(hero)

	recieveStorage := storage.InitEventStream(gameWorld.SendEvent)
	gameWorld.initStorage(recieveStorage)
	
	gameLog := &GameLog{}
	gameLog.InitLogEventStream(gameWorld.SendLog)
	fmt.Println("My Hero")
	ct.ChangeColor(ct.Cyan, true, ct.None, false)
	fmt.Println("-------")
	fmt.Print(hero.String())
	ct.ResetColor()

	commandHelp := new(CommandHelp)
	commandHelp.Init()
	ct.ChangeColor(ct.Green, true, ct.None, false)
	commandHelp.PrintHelpCommands()
	ct.ResetColor()

	pveFight := NewPveFight()
	gameWorld.addChannel(pveFight.SendEvent)

	// REPL
	ct.ChangeColor(ct.Green, true, ct.None, false)
	fmt.Print("Please enter command: ")
	ct.ResetColor()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()

		switch command {
		case "adventure":
			a := NewAdventure(hero)
			a.Embark(pveFight)
			// todo: call adventure code and pass in channel to recieve game engine messages
			// todo: don't allow user to enter new command until adventure outcome is done (wait on event?)
		case "help":
			ct.ChangeColor(ct.Green, true, ct.None, false)
			commandHelp.PrintHelpCommands()
		case "me":
			ct.ChangeColor(ct.Cyan, true, ct.None, false)
			fmt.Println("Your Hero:")
			fmt.Print(hero.String())
			fmt.Println()
			ct.ResetColor()
		case "log":
			gameLog.PrintGameLog()
			commandHelp.PrintHelpCommands()
		case "quit", "q":
			ct.ChangeColor(ct.Cyan, true, ct.None, false)
			fmt.Println()
			fmt.Println("Leaving already?? Come back soon!")
			fmt.Println()
			ct.ResetColor()
			// todo: save game state
			os.Exit(0)
		default:
			ct.ChangeColor(ct.Red, true, ct.None, false)
			fmt.Println("Unknown command.")
			ct.ResetColor()
		}

		ct.ChangeColor(ct.Green, true, ct.None, false)
		fmt.Print("Please enter command: ")
		ct.ResetColor()
	}
}

func ClearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

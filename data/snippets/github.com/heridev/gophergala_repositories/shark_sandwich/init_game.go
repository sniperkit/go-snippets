package main

import (
	"bufio"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"os"
	"strings"
)

func InitGame(ConsoleReader *bufio.Reader, storage *Storage) (*HeroSheet, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = storage.OpenRepository(dir)
	if err != nil {
		err = loadGame(ConsoleReader, storage)
		if err != nil {
			return nil, err
		}
	}

	_, err = storage.GetGameObject("shark_sandwich_game")
	if err != nil {
		err = loadGame(ConsoleReader, storage)
		if err != nil {
			return nil, err
		}
	}

	hero, err := storage.GetCurrentPlayer()
	if err != nil {
		hero, err = createNewPlayer(ConsoleReader, storage)
		if err != nil {
			return nil, err
		}
	}

	return hero, nil
}

func loadGame(ConsoleReader *bufio.Reader, storage *Storage) error {
	ct.ChangeColor(ct.Red, true, ct.None, false)
	fmt.Print("Hmm. I don't see a game in this folder. Please enter a folder location to start the game from: ")
	ct.ResetColor()
	folderPath, err := ConsoleReader.ReadString('\n')
	if err != nil {
		return err
	}

	folderPath = strings.TrimSpace(folderPath)
	err = storage.OpenRepository(folderPath)
	if err != nil {
		ct.ChangeColor(ct.Magenta, true, ct.None, false)
		fmt.Print("There is not a game there either. Let's pull down a new game. Please enter a URL to load a game: ")
		ct.ResetColor()
		remoteUrl, err := ConsoleReader.ReadString('\n')
		if err != nil {
			return err
		}

		remoteUrl = strings.TrimSpace(remoteUrl)
		err = storage.CloneRepository(remoteUrl, folderPath)
		if err != nil {
			return err
		}
		ct.ChangeColor(ct.Green, true, ct.None, false)
		fmt.Println("The game is downloaded and ready to go! Let's play!")
		fmt.Println("")
		ct.ResetColor()
	}

	ct.ChangeColor(ct.Green, true, ct.None, false)
	fmt.Println("Game found! Let's play!")
	fmt.Println("")
	ct.ResetColor()

	return nil
}

func createNewPlayer(ConsoleReader *bufio.Reader, storage *Storage) (*HeroSheet, error) {
	ct.ChangeColor(ct.Cyan, true, ct.None, false)
	fmt.Print("Looks like you're new. Tell us about your hero so you can get started. What's your name? ")
	ct.ResetColor()
	heroName, err := ConsoleReader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	playerId := strings.TrimSpace(heroName)
	hero := NewHero(playerId)
	err = storage.StorePlayer(*hero)
	if err != nil {
		return nil, err
	}

	err = storage.SetCurrentPlayer(playerId)
	if err != nil {
		return nil, err
	}

	ct.ChangeColor(ct.Green, true, ct.None, false)
	fmt.Println("That's it! You're ready to go on an adventure!")
	ct.ResetColor()
	if err != nil {
		return nil, err
	}

	return hero, nil
}

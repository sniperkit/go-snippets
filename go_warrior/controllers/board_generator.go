package controllers

import (
	"encoding/json"
	"github.com/go_warrior/abilities"
	"github.com/go_warrior/characters"
	"github.com/go_warrior/environment"
	"github.com/go_warrior/game"
	"io/ioutil"
	"os"
)

type BoardGenerator struct {
}

func NewBoardGenerator() *BoardGenerator {
	return &BoardGenerator{}
}

func (this *BoardGenerator) Generate(levelNumber string) (*characters.Warrior, *game.Board, error) {

	config, err := getConfig(levelNumber)
	if err != nil {
		return nil, nil, err
	}

	board := getBoard(config)

	warrior := getWarrior(config)

	elementMap := getElements(config, warrior)

	for y := 0; y < board.Height; y++ {

		for x := 0; x < board.Width; x++ {
			key := GenerateKey(x, y)

			space := game.NewSpace(board, x, y)

			if element, exists := elementMap[key]; exists {
				element.SetSpace(space)
				space.Element = element
			}

			board.Spaces[key] = space
		}

	}

	return warrior, board, nil
}

func getConfig(levelNumber string) (*LevelConfig, error) {
	file, err := os.Open("levels/level" + levelNumber + ".json")
	if err != nil {
		return nil, err
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	level := &LevelConfig{}

	err = json.Unmarshal(fileBytes, level)
	if err != nil {
		return nil, err
	}

	return level, nil
}

func getBoard(config *LevelConfig) *game.Board {
	return game.NewBoard(config.Board.Width, config.Board.Height)
}

func getWarrior(config *LevelConfig) *characters.Warrior {
	warrior := characters.NewWarrior()

	for _, abilityStr := range config.WarriorAbilities {
		switch abilityStr {

		case "feel":
			warrior.Abilities.Map["feel"] = &abilities.Feel{}
			break

		case "attack":
			warrior.Abilities.Map["attack"] = &abilities.Attack{}
			break

		}
	}

	return warrior
}

func getElements(config *LevelConfig, warrior *characters.Warrior) map[string]game.Element {
	elements := map[string]game.Element{}

	for _, element := range config.Elements {
		key := GenerateKey(element.X, element.Y)

		switch element.Type {
		case "warrior":
			elements[key] = warrior
			break

		case "stairs":
			elements[key] = &environment.Stairs{}
			break

		case "slug":
			elements[key] = &characters.Slug{}
			break
		}
	}

	return elements
}

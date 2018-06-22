package controllers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestLevel(t *testing.T) {
	file, err := os.Open("../levels/level00.json")
	if err != nil {
		t.Error("An error has ocurred:", err)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error("An error has ocurred:", err)
		return
	}

	config := &LevelConfig{}

	err = json.Unmarshal(fileBytes, config)
	if err != nil {
		t.Error("An error has ocurred:", err)
		return
	}

	if len(config.WarriorAbilities) < 1 {
		t.Error("No habilities were added")
		return
	}

	if config.WarriorAbilities[0] != "feel" {
		t.Error("Incorrect habilitie")
		return
	}

	if len(config.Elements) < 3 {
		t.Error("Elements missing")
		return
	}

	if config.Board.Height != 3 {
		t.Error("Invalid Board Height")
		return
	}
}

package controllers

type LevelConfig struct {
	Board            BoardData     `json:"board"`
	WarriorAbilities []string      `json:"abilities"`
	Elements         []ElementData `json:"elements"`
}

type ElementData struct {
	Type string `json:"type"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type BoardData struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

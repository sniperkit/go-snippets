package main

import (
	"math/rand"
	"strconv"
	"time"
)

type BoxesAction struct {
}

func (action BoxesAction) GetLayoutName() string {
	return "boxes_layout"
}

func (action BoxesAction) GetButtonTitle() string {
	return "Game of Boxes"
}

func (action BoxesAction) Run() {
	boxesScenario := Scenario{
		Title:       "",
		Description: "Where is a ball?",
		Answers:     Scenarios{},
	}

	//game.ClearViews()

	rand.Seed(time.Now().Unix())
	win := rand.Intn(2) + 1

	for i := 1; i <= 3; i++ {
		scenario := Scenario{
			Title: strconv.Itoa(i),
			PreDraw: func(scenario *Scenario) bool {
				game.ClearViews()
				return true
			},
			Answers: Scenarios{
				&Scenario{
					Title: "Exit",
					PreDraw: func(scenario *Scenario) bool {
						game.ClearViews()
						game.RestoreMainLayout()
						game.SwitchLocation()
						//game.SwitchLayout()
						//game.SwitchLocation()
						return false
					},
				},
			},
		}

		if i != win {
			scenario.Description = "Wrong."
		} else {
			scenario.Description = "YOU WIN!!1"
		}

		boxesScenario.Answers = append(boxesScenario.Answers, &scenario)
	}

	boxesScenario.Draw()
}

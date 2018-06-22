package main

import "github.com/seletskiy/go-android-rpc/android"

type Woman1Subgame struct {
	Subgame
	GlobalLocations map[string]Location
}

func (subgame Woman1Subgame) GetButtonTitle() string {
	return "Woman1 Title"
}

func (subgame Woman1Subgame) GetLayoutName() string {
	return "main_layout"
}

func (subgame Woman1Subgame) Enter(state *State) {
	defer android.PanicHandler()

	subgame.Subgame.Enter(state)

	locations := map[string]Location{
		"1": &BaseLocation{
			Description: `<p><b>It is pitch black.</b></p>`,
		},

		"2": &BaseLocation{
			ButtonTitle: "Try to move",
			Description: `<p>You feel different. Time goes different.</p><p>You feel your mind wandering.</p>`,
		},

		"3": &BaseLocation{
			ButtonTitle: `"Who are you?"`,
			Description: `<p>"You? Things? Cold. Solid. Energy. Dark."</p>`,
		},

		"4.1": &BaseLocation{
			ButtonTitle: "Play.",
			Description: "yaaaaa",
		},

		"4.2": &BaseLocation{
			ButtonTitle: "Toys.",
			Description: "huyaaaa",
		},

		"4.3": &BaseLocation{
			ButtonTitle: "Mother.",
			Description: "mama-ama-criminal",
		},

		"4.4": &BaseLocation{
			ButtonTitle: "Father.",
			Description: "where is mother?",
		},

		"5": &BaseLocation{
			ButtonTitle: "Run away",
			Description: "<p>It is pitch black. There is nowhere to run.</p>",
		},

		"6": &BaseLocation{
			ButtonTitle: `End`,
			Description: `<p>Father takes you</p>`,
		},
	}

	locations["1"].Link(locations["2"])

	locations["2"].Link(locations["3"])

	locations["3"].Link(locations["4.1"])
	locations["3"].Link(locations["4.2"])
	locations["3"].Link(locations["4.3"])
	locations["3"].Link(locations["4.4"])

	locations["4.1"].Link(locations["5"])
	locations["4.2"].Link(locations["5"])
	locations["4.3"].Link(locations["5"])
	locations["4.4"].Link(locations["5"])

	locations["5"].Link(locations["6"])

	locations["6"].Link(subgame.GlobalLocations["bunk"])

	subgame.SetLayoutName("main_layout")
	subgame.SetLocation(locations["1"])

	subgame.Start()
}

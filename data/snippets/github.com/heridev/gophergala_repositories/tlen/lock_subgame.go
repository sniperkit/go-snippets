package main

import (
	"log"
	"strconv"

	"github.com/seletskiy/go-android-rpc/android"
	"github.com/seletskiy/go-android-rpc/android/sdk"
)

type LockedSubgame struct {
	Subgame
	GlobalLocations map[string]Location
}

func (subgame LockedSubgame) GetButtonTitle() string {
	return "Captain cabinet (digital lock)"
}

func (subgame LockedSubgame) GetLayoutName() string {
	return "locked_layout"
}

func (subgame LockedSubgame) Enter(state *State) {
	defer android.PanicHandler()

	subgame.Subgame.Enter(state)

	locations := map[string]Location{
		"1": &AwaitLockedLocation{
			BaseLocation: BaseLocation{
				Description: "DESC",
			},
			Subgame: &subgame,
		},
	}

	subgame.SetLayoutName("locked_layout")
	subgame.SwitchLayout()
	subgame.SetLocation(locations["1"])

	subgame.Start()
}

type AwaitLockedLocation struct {
	BaseLocation
	Relations map[int][]int
	Buttons   map[int]sdk.ToggleButton
	States    map[int]bool
	Subgame   *LockedSubgame
}

func (location AwaitLockedLocation) Enter(state *State) {
	defer android.PanicHandler()
	// buttons:
	//     1 2 3
	//     4 5 6
	//     7 8 9
	location.Relations = make(map[int][]int)
	location.Buttons = make(map[int]sdk.ToggleButton)
	location.States = make(map[int]bool)

	location.States = map[int]bool{
		1: false, 2: false, 3: false,
		4: false, 5: false, 6: false,
		7: false, 8: false, 9: false,
	}

	location.Relations = map[int][]int{
		1: []int{
			2, 3, 4, 7,
		},
		2: []int{
			1, 3, 5, 8,
		},
		3: []int{
			1, 2, 6, 9,
		},
		4: []int{
			1, 7, 5, 6,
		},
		5: []int{
			2, 8, 4, 6,
		},
		6: []int{
			4, 5, 3, 9,
		},
		7: []int{
			1, 4, 8, 9,
		},
		8: []int{
			2, 5, 7, 9,
		},
		9: []int{
			3, 6, 7, 8,
		},
	}

	for i := 1; i <= 9; i++ {
		button := android.GetViewById(
			"toggle_button_" + strconv.Itoa(i)).(sdk.ToggleButton)

		android.OnClick(button,
			AwaitLockedLocationHandler{
				&location, i})

		location.Buttons[i] = button
	}

	log.Printf("%#v", location.Relations)
}

type AwaitLockedLocationHandler struct {
	location *AwaitLockedLocation
	index    int
}

func (handler AwaitLockedLocationHandler) OnClick() {
	handler.location.OnClick(handler.index)
}

func (location AwaitLockedLocation) OnClick(index int) {
	defer android.PanicHandler()
	location.States[index] = !location.States[index]

	for _, otherIndex := range location.Relations[index] {
		state := !location.States[otherIndex]
		location.Buttons[otherIndex].SetChecked(state)
		location.States[otherIndex] = state
	}

	log.Printf("%#v", location.States)
	for _, state := range location.States {
		if !state {
			return
		}
	}

	// sorry mom, I'm used java
	location.Subgame.SetLayoutName("main_layout")
	location.Subgame.SwitchLayout()

	location.Subgame.SetLocation(location.Subgame.GlobalLocations["bunk"])
	log.Printf("%#v", location.Subgame)
	location.Subgame.SwitchLocation()
}

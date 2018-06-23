package main

import (
	"github.com/seletskiy/go-android-rpc/android"
	"github.com/seletskiy/go-android-rpc/android/sdk"
)

type CaressCatSubgame struct {
	Subgame

	OriginLocation Location
}

func (subgame CaressCatSubgame) GetButtonTitle() string {
	return "Play with cat"
}

func (subgame CaressCatSubgame) GetLayoutName() string {
	return "main_layout"
}

func (subgame CaressCatSubgame) Enter(state *State) {
	subgame.Subgame.Enter(state)

	defer android.PanicHandler()

	main := &CaressCatLocation{game: &subgame}
	main.Link(subgame.OriginLocation)

	subgame.SetLocation(main)
	subgame.Start()
}

type CaressCatOnClickHandler struct {
	button sdk.Button
	state  *State
}

func (handler CaressCatOnClickHandler) OnClick() {
	handler.button.PerformHapticFeedback(0)
	handler.state.Progress++
	handler.state.Cat++
	handler.state.MoveCounter++
}

type CaressCatLocation struct {
	BaseLocation

	game *CaressCatSubgame
}

func (location CaressCatLocation) GetDescription() string {
	return "<p>" + Cat + " is looking at you.</p>"
}

func (location CaressCatLocation) GetHeader() string {
	return "<p>You sit down near the cat.</p>"
}

func (location CaressCatLocation) Enter(state *State) {
	caressButton := location.game.CreateView(
		"android.widget.Button",
	).(sdk.Button)

	android.SetTextFromHtml(caressButton, "Caress "+Cat)

	android.OnClick(caressButton, CaressCatOnClickHandler{
		caressButton,
		state,
	})

	location.game.AttachView(caressButton.View)
}

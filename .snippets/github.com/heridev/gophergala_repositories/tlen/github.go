package main

import "github.com/seletskiy/go-android-rpc/android"

type GithubSubgame struct {
	Subgame
	GlobalLocations map[string]Location
}

func (subgame GithubSubgame) GetButtonTitle() string {
	return "Visit our github repo"
}

func (subgame GithubSubgame) GetLayoutName() string {
	return "main_layout"
}

func (subgame GithubSubgame) Enter(state *State) {
	defer android.PanicHandler()

	subgame.Subgame.Enter(state)

	android.OpenWebPage("http://github.com/seletskiy/tlen")

	//subgame.SetLayoutName("main_layout")
	//subgame.SetLocation(&BaseLocation{
		//ButtonTitle:

	subgame.SetLocation(subgame.GlobalLocations["lobby"])
	subgame.Start()
}

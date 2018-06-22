package main

import "github.com/seletskiy/go-android-rpc/android"

type Stage1ProxyLocation struct {
	Subgame

	OldWay Location
	NewWay Location
}

func (subgame Stage1ProxyLocation) GetButtonTitle() string {
	return subgame.OldWay.GetButtonTitle()
}

func (subgame Stage1ProxyLocation) GetDescription() string {
	return subgame.OldWay.GetDescription()
}

func (subgame Stage1ProxyLocation) Enter(state *State) {
	defer android.PanicHandler()

	subgame.Subgame.Enter(state)

	if subgame.IsNewWayAvailable() {
		subgame.SetLocation(subgame.NewWay)
	} else {
		subgame.SetLocation(subgame.OldWay)
	}

	subgame.Start()
}

func (subgame Stage1ProxyLocation) IsNewWayAvailable() bool {
	return true
}

package main

type Subgame struct {
	BaseLocation
	Game
}

func (subgame *Subgame) Enter(state *State) {
	subgame.state = state
}

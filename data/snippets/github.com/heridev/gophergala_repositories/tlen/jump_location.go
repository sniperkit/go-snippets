package main

type JumpLocation struct {
	Subgame
	BaseLocation

	Jump Location
}

func (location JumpLocation) Enter(state *State) {
	location.Subgame.Enter(state)

	location.SetLocation(location.Jump)

	location.Start()
}

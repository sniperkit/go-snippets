package main

import (
	"github.com/gophergala/GopherKombat/common/game"
	"log"
)

type Direction int

const (
	N  Direction = 0
	NE Direction = 1
	E  Direction = 2
	SE Direction = 3
	S  Direction = 4
	SW Direction = 5
	W  Direction = 6
	NW Direction = 7
)

var directions = map[Direction]struct {
	dx int
	dy int
}{
	N:  {dx: 0, dy: -1},
	NE: {dx: 1, dy: -1},
	E:  {dx: 1, dy: 0},
	SE: {dx: 1, dy: 1},
	S:  {dx: 0, dy: 1},
	SW: {dx: -1, dy: 1},
	W:  {dx: -1, dy: 0},
	NW: {dx: -1, dy: -1},
}

type UserData struct {
}

type Gopher struct {
	Id    int
	state *game.State
	// Custom user data
	Data UserData
}

func (gopher *Gopher) Update(state *game.State) {
	log.Printf("Turn: %#v\n", state)
	gopher.state = state
}

func (gopher *Gopher) Health() int {
	return gopher.state.Health
}

func (gopher *Gopher) Ammo() int {
	return gopher.state.Ammo
}

func (gopher *Gopher) X() int {
	return gopher.state.Me.X
}

func (gopher *Gopher) Y() int {
	return gopher.state.Me.Y
}

func (gopher *Gopher) Nearby() []game.GopherData {
	return gopher.state.Nearby
}

func (gopher *Gopher) Nop() *game.Action {
	return &game.Action{Code: game.Nop}
}

func (gopher *Gopher) Move(dir Direction) *game.Action {
	return &game.Action{Code: game.Move,
		X: gopher.X() + directions[dir].dx,
		Y: gopher.Y() + directions[dir].dy}
}

func (gopher *Gopher) Shoot(target *game.GopherData) *game.Action {
	return &game.Action{Code: game.Shoot,
		X: target.X,
		Y: target.Y}
}

func (gopher *Gopher) Init() {
	log.Printf("Init gopher %d\n", gopher.Id)
}

func (gopher *Gopher) Turn() *game.Action {
	return gopher.Nop()
}

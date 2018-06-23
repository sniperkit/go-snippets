package main

import (
	"fmt"
)

type User struct {
	Id             int
	Name, Location string
}

type Player struct {
	*User
	GameId int
}

func (u *Player) Greetings() string {
	return fmt.Sprintf("Hi %s from %s with Game ID: %d", u.Name, u.Location, u.GameId)
}

func newPlayer(id int, name string, location string, gameId int) *Player {
	return &Player{
		User:   &User{Id: id, Name: name, Location: location},
		GameId: gameId,
	}
}

func main() {
	player := newPlayer(22, "Nexus", "TN", 145)
	fmt.Println(player.Greetings())
}

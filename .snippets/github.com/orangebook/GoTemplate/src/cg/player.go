package cg

import "fmt"

type Player struct {
	Name  string
	Level int
	Exp   int
	Room  int
	mq    chan *Message
}

func NewPlayer() *Player {
	m := make(chan *Message, 1024)
	player := &Player{"", 0, 0, 0, m}
	go func(p *Player) {
		msg := <-p.mq
		fmt.Println(p.Name, "received message:", msg.Content)
	}(player)
	return player
}
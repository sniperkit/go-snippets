package main

import (
	"fmt"
	"strings"
)

type (
	FightEvent struct {
		Won       bool
		EnemyName string
		Event
	}

	GameWorld struct {
		Hero      *HeroSheet
		SendEvent chan string
		SendLog   chan LogEvent
	}
)

func NewGameWorld(hero *HeroSheet) *GameWorld {
	return &GameWorld{hero, make(chan string, 10), make(chan LogEvent, 100)}
}

func (g *GameWorld) initStorage(events chan string) {
	for event := range events {
		if strings.HasPrefix(event, "You Won") {
			g.Hero.Xp = g.Hero.Xp + 10
			log := LogEvent{
				"Fight: " + event,
				int(g.Hero.Xp),
				int(g.Hero.Life),
				int(g.Hero.Speed),
				int(g.Hero.Power),
				int(g.Hero.Ancestry),
			}
			g.SendLog <- log
		} else {
			log := LogEvent{
				"Fight: " + event,
				int(g.Hero.Xp),
				int(g.Hero.Life),
				int(g.Hero.Speed),
				int(g.Hero.Power),
				int(g.Hero.Ancestry),
			}
			g.SendLog <- log
		}
	}
}

func (g *GameWorld) addChannel(events chan interface{}) {
	go func() {
		for {
			e := <-events
			switch event := e.(type) {
			case FightEvent:
				message := fmt.Sprintf("Fight: %s\n", event.String())
				g.SendEvent <- event.String()
				if event.Won {
					g.Hero.Xp = g.Hero.Xp + 10
				}
				log := LogEvent{
					message,
					int(g.Hero.Xp),
					int(g.Hero.Life),
					int(g.Hero.Speed),
					int(g.Hero.Power),
					int(g.Hero.Ancestry),
				}
				g.SendLog <- log
			}
		}
	}()
}

func (f *FightEvent) String() string {
	if f.Won {
		return "You Won a fight with a " + f.EnemyName
	} else {
		return "You Lost a fight with a " + f.EnemyName
	}
}

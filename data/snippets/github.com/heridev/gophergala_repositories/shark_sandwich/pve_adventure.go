package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"math/rand"
)

const (
	ADVENTURE_TYPE_ENCOUNTER AdventureType = 1
	ADVENTURE_TYPE_WANDER    AdventureType = 2
	ADVENTURE_TYPE_DISCOVERY AdventureType = 3
)

type (
	AdventureType int
	Adventure     struct {
		Type AdventureType
		*HeroSheet
	}
)

func NewAdventure(h *HeroSheet) *Adventure {
	return &Adventure{
		Type:      generateAdventure(),
		HeroSheet: h,
	}
}

func generateAdventure() AdventureType {
	return AdventureType(random(1, 3))
}

func generateEnemy() string {
	names := []string{
		"Snorlax",
		"Bear",
		"Goblin",
		"Hill Giant",
		"Shark out of water",
	}

	return names[rand.Intn(len(names))]
}

func generateNothing() string {
	messages := []string{
		"You wandered around a bit, got bored, and went home.",
		"You forgot to go to the bathroom before you left, and went home.",
		"You got turned around somewhere and ended up at home.",
		"You wandered around, saw a stump, a leaf, and a caterpillar, then went home.",
		"You had an anticlimactic experience and wound up back at home with a cup of tea.",
	}

	return messages[rand.Intn(len(messages))]
}

func (a *Adventure) Embark(pve *PveFight) {
	ct.ChangeColor(ct.Magenta, true, ct.None, false)
	switch a.Type {
	case ADVENTURE_TYPE_DISCOVERY:
		fmt.Println("You didn't discover anything, too bad.")
	case ADVENTURE_TYPE_ENCOUNTER:
		enemyName := generateEnemy()
		fmt.Printf("A wild %s appeared.\n", enemyName)
		enemey := NewEnemy(a.HeroSheet, enemyName)

		if enemey.Life > a.HeroSheet.BaseStats.Life {
			fmt.Println(" It's a tough one!")
		}
		if enemey.Power > a.HeroSheet.BaseStats.Power {
			fmt.Println(" It hits pretty hard...")
		}
		if enemey.Speed > a.HeroSheet.BaseStats.Speed {
			fmt.Printf(" Faster than your average %s.\n", enemyName)
		}
		fmt.Printf("%s attacks you! Check 'log' to see the result.\n", enemyName)
		pve.Fight(a.HeroSheet, enemey)
	case ADVENTURE_TYPE_WANDER:
		fmt.Println(generateNothing())
	}
	ct.ResetColor()
}

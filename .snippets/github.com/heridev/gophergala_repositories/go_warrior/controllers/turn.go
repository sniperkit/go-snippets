package controllers

import (
	//"fmt"
	"github.com/go_warrior/abilities"
	"github.com/go_warrior/game"
)

type UserFunction func()

type TurnController struct {
	Board *game.Board
}

func (this *TurnController) Start() {

}

func (this *TurnController) getAllCharacters() (map[int][]abilities.Performer, int) {
	performersByInitiative := map[int][]abilities.Performer{}
	maxInitiative := 0

	for _, space := range this.Board.Spaces {
		if space.Element != nil {
			if space.Element.GetType() != "stairs" {
				getByInitiative(space, &performersByInitiative, &maxInitiative)
			}
		}
	}

	return performersByInitiative, maxInitiative
}

func (this *TurnController) resolveTurns(f UserFunction, performers map[int][]abilities.Performer, maxInitiative int) {
	for i := maxInitiative; i >= 1; i-- {
		initiativeLayer := performers[i]
		for _, performer := range initiativeLayer {
			if performer.GetType() == "warrior" {
				f()
			}
		}
	}
}

func getByInitiative(space *game.Space, performerMap *map[int][]abilities.Performer, maxInitiative *int) {
	performer := space.Element.(abilities.Performer)
	initiative := performer.GetInitiative()

	if initiative > *maxInitiative {
		*maxInitiative = initiative
	}

	pmap := *performerMap

	if _, exists := pmap[initiative]; !exists {
		pmap[initiative] = []abilities.Performer{}
	}

	pmap[initiative] = append(pmap[initiative], performer)
}

func isLevelSucceded() bool {
	return game.STATE == game.SUCCESS
}

package abilities

import (
	"errors"
	"github.com/go_warrior/game"
)

type Ability interface {
	Perform(Performer, game.Direction) interface{}
}

type Abilities struct {
	Performer Performer
	Map       map[string]Ability
}

func (this *Abilities) Feel(direction game.Direction) *game.Space {
	if ability, exists := this.Map["feel"]; !exists {

		panic(errors.New("No such ability"))

	} else {

		return ability.Perform(this.Performer, direction).(*game.Space)

	}

	return nil
}

func (this *Abilities) Attack(direction game.Direction) {
	if ability, exists := this.Map["attack"]; !exists {

		panic(errors.New("No such ability"))

	} else {

		ability.Perform(this.Performer, direction)

	}
}

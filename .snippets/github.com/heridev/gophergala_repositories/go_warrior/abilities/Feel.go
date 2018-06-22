package abilities

import "github.com/go_warrior/game"

type Feel struct{}

func (this *Feel) Perform(performer Performer, direction game.Direction) interface{} {
	space := performer.GetSpace()

	return space.GetNext(direction)
}

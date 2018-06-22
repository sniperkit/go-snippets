package abilities

import (
	"fmt"
	"github.com/go_warrior/game"
)

type Attack struct{}

func (this *Attack) Perform(performer Performer, direction game.Direction) interface{} {
	space := performer.GetSpace()

	nextSpace := space.GetNext(direction)

	if nextSpace == nil {
		return nil
	}

	target := nextSpace.Element

	if target == nil {
		return nil
	}

	if target.GetType() != "stairs" {
		fmt.Printf("%s attacked %s", performer.GetType(), target.GetType())
	}

	return nil
}

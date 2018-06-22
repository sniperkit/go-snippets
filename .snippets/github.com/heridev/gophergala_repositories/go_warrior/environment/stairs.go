package environment

import "github.com/go_warrior/game"

type Stairs struct {
}

func (this *Stairs) SetSpace(space *game.Space) {
}

func (this *Stairs) GetSpace() *game.Space {
	return nil
}

func (this *Stairs) GetSprite() string {
	return ">"
}

func (this *Stairs) GetType() string {
	return "stairs"
}

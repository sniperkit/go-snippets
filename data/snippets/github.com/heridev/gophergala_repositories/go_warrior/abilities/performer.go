package abilities

import "github.com/go_warrior/game"

type Performer interface {
	GetSpace() *game.Space
	GetInitiative() int
	GetType() string
}

package characters

import "github.com/go_warrior/game"

type Slug struct {
	Health       int
	AttackPoints int
	Space        *game.Space
}

func GetSlug(space, x, y int) *Warrior {
	return &Warrior{
		Health:       10,
		AttackPoints: 3,
	}
}

func (this *Slug) SetSpace(space *game.Space) {
	this.Space = space
}

func (this *Slug) GetSpace() *game.Space {
	return this.Space
}

func (this *Slug) GetSprite() string {
	return "s"
}

func (this *Slug) GetType() string {
	return "slug"
}

func (this *Slug) GetInitiative() int {
	return 1
}

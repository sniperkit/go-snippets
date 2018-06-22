package game

type Element interface {
	SetSpace(*Space)
	GetSpace() *Space
	GetSprite() string
	GetType() string
}

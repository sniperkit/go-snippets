package game

type Board struct {
	Width  int
	Height int
	Spaces map[string]*Space
}

func NewBoard(width, height int) *Board {
	return &Board{
		Width:  width,
		Height: height,
		Spaces: map[string]*Space{},
	}
}

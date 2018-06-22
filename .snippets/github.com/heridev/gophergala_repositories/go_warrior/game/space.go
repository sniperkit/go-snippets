package game

import "strconv"

type Space struct {
	Element Element
	Board   *Board
	X       int
	Y       int
}

func NewSpace(board *Board, x, y int) *Space {
	return &Space{
		X:     x,
		Y:     y,
		Board: board,
	}
}

func (this *Space) Empty() bool {
	return this.Element == nil
}

func (this *Space) Enemy() bool {
	return !this.Empty() && this.Element.GetType() != "stairs"
}

func (this *Space) Stairs() bool {
	return !this.Empty() && this.Element.GetType() == "stairs"
}

func (this *Space) GetNext(direction Direction) *Space {
	key := strconv.Itoa(this.X+direction.X) + "-" + strconv.Itoa(this.Y+direction.Y)
	if space, exists := this.Board.Spaces[key]; exists {
		return space
	}
	return nil
}

package board

import "github.com/gophergala/gallivanting_gophers/pieces"

// DirLeft and friends are directional indicators (bitflags)
const (
	DirLeft   = 1
	DirTop    = 2
	DirRight  = 4
	DirBottom = 8
)

// FixBoard is the second implementation of a board (this one actually works! :P)
// It seems that it is overly complicated (utilizing bitsets for storage) but
// this is done purposefully to speed up the board analysis (calculations against
// generated boards). This is based on some benchmarking I've done earlier in
// another project.
type FixBoard struct {
	Spaces  [16]uint         `json:"tiles"`
	Gophers []*pieces.Gopher `json:"gophers"`
	Goals   []pieces.Goal    `json:"goals"`

	locations []int
}

// Creates a new empty board.
func newEmptyBoard() *FixBoard {
	b := &FixBoard{
		Gophers: make([]*pieces.Gopher, 0),
		Goals:   make([]pieces.Goal, 0),
	}

	return b
}

// SetWall sets a specific space's wall based on the direction from the space.
func (b *FixBoard) SetWall(s int, w int) {
	b.modifySpace(s, w)
}

// SetWalls sets all of the walls associated with a specific space.
func (b *FixBoard) SetWalls(s int, left bool, top bool, right bool, bottom bool) {
	modifier := 0

	if left {
		modifier |= DirLeft
	}

	if top {
		modifier |= DirTop
	}

	if right {
		modifier |= DirRight
	}

	if bottom {
		modifier |= DirBottom
	}

	b.modifySpace(s, modifier)
}

// modifies the space based on
func (b *FixBoard) modifySpace(s int, w int) {
	row := s / 16
	column := s % 16

	// Left - Only if not in leftmost column
	if column > 0 && (w&DirLeft) > 0 {
		b.setBit(row, uint(31-(column-1)))
	}

	// Top - Only if not in topmost row
	if row > 0 && (w&DirTop) > 0 {
		b.setBit(row-1, uint(15-column))
	}

	// Right - Only if not in rightmost column
	if column < 15 && (w&DirRight) > 0 {
		b.setBit(row, uint(31-column))
	}

	// Bottom - Only if not in bottom most row
	if row < 15 && (w&DirBottom) > 0 {
		b.setBit(row, uint(15-column))
	}
}

// Sets a bit based on offset and row.
func (b *FixBoard) setBit(i int, l uint) {
	if i < 0 {
		return
	}
	b.Spaces[i] |= 1 << l
}

// Checks a bit based on offset and row.
func (b *FixBoard) hasBit(i int, l uint) bool {
	if i < 0 {
		return true
	}
	return b.Spaces[i]&(1<<l) > 0
}

// HasWall checks for a wall to exist based on space index and direction.
func (b *FixBoard) HasWall(s int, w int) bool {
	row := s / 16
	column := s % 16

	if w == DirLeft {
		return b.hasBit(row, uint(31-(column-1))) || column == 0
	}

	if w == DirTop {
		return b.hasBit(row-1, uint(15-column)) || row == 0
	}

	if w == DirRight {
		return b.hasBit(row, uint(31-column)) || column == 15
	}

	if w == DirBottom {
		return b.hasBit(row, uint(15-column)) || row == 15
	}

	return false
}

// MoveGopher moves a gopher by type based on direction.
func (b *FixBoard) MoveGopher(g byte, d int) {
	for _, v := range b.Gophers {
		if v.Type != g {
			continue
		}

		var smod int
		s := v.Location

		if d == DirLeft {
			smod = -1
		} else if d == DirTop {
			smod = -16
		} else if d == DirRight {
			smod = 1
		} else if d == DirBottom {
			smod = 16
		}

		// canMove := !b.HasWall(s, d)
		canMove := b.canMove(s, d, false)

		for canMove {
			s = s + smod
			// canMove = !b.HasWall(s, d)
			canMove = b.canMove(s, d, false)
		}

		v.Location = s
	}
}

// canMove checks to see if a move is valid. The 'i' parameter is for
// ignoring gopher locations when calculating (this is used by solver
// and should be changed to support a cleaner solving mechanism).
func (b *FixBoard) canMove(s int, d int, i bool) bool {
	var smod int

	if d == DirLeft {
		smod = -1
	} else if d == DirTop {
		smod = -16
	} else if d == DirRight {
		smod = 1
	} else if d == DirBottom {
		smod = 16
	}

	proposed := s + smod

	if !i {
		for _, v := range b.Gophers {
			if v.Location == proposed {
				return false
			}
		}
	}

	return !b.HasWall(s, d)
}

// MoveCheck allows you to determine the new location without actually moving the piece
// You can provide a list of locations to excluse as well (temporary gophers)
func (b *FixBoard) MoveCheck(s int, d int, l []int) int {
	var smod int

	result := s

	if d == DirLeft {
		smod = -1
	} else if d == DirTop {
		smod = -16
	} else if d == DirRight {
		smod = 1
	} else if d == DirBottom {
		smod = 16
	}

	canMove := b.canMove(s, d, true)
	proposed := s + smod
	for _, v := range l {
		if proposed == v {
			canMove = false
		}
	}

	for canMove {
		s = s + smod
		canMove = b.canMove(s, d, true)
		proposed = s + smod
		for _, v := range l {
			if proposed == v {
				canMove = false
			}
		}
	}

	return result
}

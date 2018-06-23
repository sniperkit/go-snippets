package board

// Solver where did all my time go? :(
type Solver struct {
	Board           *FixBoard
	Solutions       []*Solution
	cachedPositions map[uint64]struct{}
}

type Solution struct {
	Gopher byte
	Number int
	Moves  []Move
}

type Move struct {
	Gopher    byte
	Direction int
}

// NewSolver ...
func NewSolver(b *FixBoard) *Solver {
	s := &Solver{
		Board:           b,
		Solutions:       make([]*Solution, 0),
		cachedPositions: make(map[uint64]struct{}),
	}

	return s
}

// Solve ...
func (s *Solver) Solve() {

}

func (s *Solver) possibleMoves(position uint64) *[]Move {
	r := make([]Move, 0)

	return &r
}

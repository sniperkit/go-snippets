package game

type Direction struct {
	X int
	Y int
}

var (
	Right    = Direction{X: 1, Y: 0}
	Left     = Direction{X: -1, Y: 0}
	Backward = Direction{X: 0, Y: 1}
	Forward  = Direction{X: 0, Y: -1}
)

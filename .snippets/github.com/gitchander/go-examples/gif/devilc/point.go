package devilc

type Point struct {
	X, Y float64
}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

func (p Point) MulScalar(s float64) Point {
	return Point{p.X * s, p.Y * s}
}

func (p Point) DivScalar(s float64) Point {
	s = 1 / s
	return Point{p.X * s, p.Y * s}
}

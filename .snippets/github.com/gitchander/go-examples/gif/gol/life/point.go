package life

type Point struct {
	X, Y int
}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

func (p Point) Equal(q Point) bool {
	if p.X != q.X {
		return false
	}
	if p.Y != q.Y {
		return false
	}
	return true
}

// Neighbors shifts
var neShifts = []Point{
	Point{-1, -1},
	Point{0, -1},
	Point{+1, -1},
	Point{-1, 0},
	Point{+1, 0},
	Point{-1, +1},
	Point{0, +1},
	Point{+1, +1},
}

type Points []Point

// Contained
func (ps Points) Include(p Point) bool {
	for _, q := range ps {
		if q.Equal(p) {
			return true
		}
	}
	return false
}

func (ps Points) Exclude(p Point) bool {
	return !ps.Include(p)
}

func (ps Points) Neighbors(p Point) (count int) {
	for _, ns := range neShifts {
		n := p.Add(ns)
		for _, q := range ps {
			if q.Equal(n) {
				count++
			}
		}
	}
	return
}

func (ps Points) Move(p Point) {
	for i := range ps {
		ps[i] = ps[i].Add(p)
	}
}

func (ps Points) Clone() Points {
	qs := make([]Point, len(ps))
	copy(qs, ps)
	return qs
}

func (a Points) Equal(b Points) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[Point]struct{})
	for _, p := range a {
		m[p] = struct{}{}
	}
	for _, p := range b {
		if _, ok := m[p]; !ok {
			return false
		}
	}
	return true
}

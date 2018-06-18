package life

type Nexter interface {
	Next(curr []Point) (next []Point)
}

type Nexter1 struct{}

var _ Nexter = Nexter1{}

func (Nexter1) Next(curr []Point) (next []Point) {

	// dead
	for _, c := range curr {
		for _, ns := range neShifts {
			p := c.Add(ns)
			if Points(next).Exclude(p) {
				if Points(curr).Exclude(p) {
					n := Points(curr).Neighbors(p)
					if n == 3 {
						next = append(next, p)
					}
				}
			}
		}
	}

	// alive
	for _, p := range curr {
		n := Points(curr).Neighbors(p)
		if (n == 2) || (n == 3) {
			next = append(next, p)
		}
	}

	return
}

type pointState struct {
	Point Point
	Alive bool
}

type Nexter2 struct{}

var _ Nexter = Nexter2{}

func (Nexter2) Next(curr []Point) (next []Point) {

	hash := hashPointToBytes

	m := make(map[hashBytes]*pointState)

	for _, c := range curr {
		m[hash(c)] = &pointState{c, true}
	}
	for _, c := range curr {
		for _, ns := range neShifts {
			var (
				p = c.Add(ns)
				h = hash(p)
			)
			if _, ok := m[h]; !ok {
				m[h] = &pointState{p, false}
			}
		}
	}

	for _, node := range m {
		n := 0
		for _, ns := range neShifts {
			p := node.Point.Add(ns)
			if pn, ok := m[hash(p)]; ok && pn.Alive {
				n++
			}
		}
		if node.Alive {
			if (n == 2) || (n == 3) {
				next = append(next, node.Point)
			}
		} else {
			if n == 3 {
				next = append(next, node.Point)
			}
		}
	}

	return
}

type Nexter3 struct{}

var _ Nexter = Nexter3{}

func (Nexter3) Next(curr []Point) (next []Point) {

	hash := hashPointToUint64

	m := make(map[uint64]*pointState)

	for _, c := range curr {
		m[hash(c)] = &pointState{c, true}
	}
	for _, c := range curr {
		for _, ns := range neShifts {
			var (
				p = c.Add(ns)
				h = hash(p)
			)
			if _, ok := m[h]; !ok {
				m[h] = &pointState{p, false}
			}
		}
	}

	for _, node := range m {
		n := 0
		for _, ns := range neShifts {
			p := node.Point.Add(ns)
			if pn, ok := m[hash(p)]; ok && pn.Alive {
				n++
			}
		}
		if node.Alive {
			if (n == 2) || (n == 3) {
				next = append(next, node.Point)
			}
		} else {
			if n == 3 {
				next = append(next, node.Point)
			}
		}
	}

	return
}

type Nexter4 struct{}

var _ Nexter = Nexter4{}

func (Nexter4) Next(curr []Point) (next []Point) {

	m := make(map[Point]bool)

	for _, c := range curr {
		m[c] = true
	}
	for _, c := range curr {
		for _, ns := range neShifts {
			p := c.Add(ns)
			if _, ok := m[p]; !ok {
				m[p] = false
			}
		}
	}

	for c, alive := range m {
		n := 0
		for _, ns := range neShifts {
			p := c.Add(ns)
			if alive, ok := m[p]; ok && alive {
				n++
			}
		}
		if alive {
			if (n == 2) || (n == 3) {
				next = append(next, c)
			}
		} else {
			if n == 3 {
				next = append(next, c)
			}
		}
	}

	return
}

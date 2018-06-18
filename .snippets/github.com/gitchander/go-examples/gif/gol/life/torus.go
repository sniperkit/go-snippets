package life

import (
	"github.com/gitchander/go-examples/gif/gol/bitmap"
)

type Torus struct {
	nX, nY     int
	curr, next *bitmap.Bitmap
}

func NewTorus(nX, nY int) *Torus {
	if nX < 0 {
		nX = 0
	}
	if nY < 0 {
		nY = 0
	}
	n := nX * nY
	return &Torus{
		nX:   nX,
		nY:   nY,
		curr: bitmap.New(n),
		next: bitmap.New(n),
	}
}

func (t *Torus) Size() (nX, nY int) {
	return t.nX, t.nY
}

func (t *Torus) Empty() bool {
	return (t.nX == 0) || (t.nY == 0)
}

func (t *Torus) swapBuffers() {
	t.curr, t.next = t.next, t.curr
}

func (t *Torus) Next() {
	if t.Empty() {
		return
	}
	for y := 0; y < t.nY; y++ {
		for x := 0; x < t.nX; x++ {
			var (
				i = t.offset(x, y)
				n = t.neighbors(x, y)
			)
			v, _ := t.curr.Get(i)
			if v {
				v = (n == 2) || (n == 3)
			} else {
				v = (n == 3)
			}
			t.next.Set(i, v)
		}
	}
	t.swapBuffers()
}

func (t *Torus) neighbors(x, y int) (n int) {
	for _, ns := range neShifts {
		offset := t.offset(x+ns.X, y+ns.Y)
		v, _ := t.curr.Get(offset)
		if v {
			n++
		}
	}
	return
}

func (t *Torus) offset(x, y int) int {
	x = mod(x, t.nX)
	y = mod(y, t.nY)
	return y*t.nX + x
}

func mod(x, y int) int {
	t := x % y
	if t < 0 {
		t += y
	}
	return t
}

func (t *Torus) Get(x, y int) bool {
	if t.Empty() {
		return false
	}
	v, _ := t.curr.Get(t.offset(x, y))
	return v
}

func (t *Torus) Set(x, y int, v bool) {
	if t.Empty() {
		return
	}
	t.curr.Set(t.offset(x, y), v)
}

func (t *Torus) Clone() *Torus {
	c := NewTorus(t.nX, t.nY)
	n := t.nX * t.nY
	for i := 0; i < n; i++ {
		v, _ := t.curr.Get(i)
		c.curr.Set(i, v)
	}
	return c
}

func (t1 *Torus) Equal(t2 *Torus) bool {
	if (t1.nX != t2.nX) || (t1.nY != t2.nY) {
		return false
	}
	n := t1.nX * t1.nY
	for i := 0; i < n; i++ {
		v1, _ := t1.curr.Get(i)
		v2, _ := t2.curr.Get(i)
		if v1 != v2 {
			return false
		}
	}
	return true
}

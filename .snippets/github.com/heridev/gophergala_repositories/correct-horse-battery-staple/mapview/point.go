package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type Point struct {
	js.Object
}

func NewPoint(x, y int64) *Point {
	return &Point{
		Object: L.Call("point", x, y),
	}
}

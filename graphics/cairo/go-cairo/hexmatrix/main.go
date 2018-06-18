package main

import (
	"fmt"

	"github.com/gitchander/cairo"

	"github.com/gitchander/heuristic/math/graph2d"
	"github.com/gitchander/heuristic/math/hexm"
)

func getCenter(surface *cairo.Surface) graph2d.Vector {
	return graph2d.Vector{
		float32(surface.GetWidth()) * 0.5,
		float32(surface.GetHeight()) * 0.5,
	}
}

func main() {

	mSize, _ := hexm.NewSize(3, 3, 3)
	m := hexm.NewMatrix(mSize)

	const cellRadius = float32(56.2)

	texture, status := cairo.NewSurfaceFromPNG("./hexagone-gray.png")
	if status != cairo.STATUS_SUCCESS {
		fmt.Println(status)
		return
	}

	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, 512, 512)
	defer surface.Destroy()

	c, status := cairo.NewCanvas(surface)
	if status != cairo.STATUS_SUCCESS {
		fmt.Println(cairo.StatusToString(status))
		return
	}
	defer c.Destroy()

	shift := getCenter(surface).Sub(getCenter(texture))

	for I := m.NewIterator(); !I.Done(); I.Next() {

		coord, _, _ := I.Current()

		v := hexm.CoordToVector(coord)
		v = v.MulScalar(cellRadius)
		v = v.Add(shift)

		c.SetSourceSurface(texture, float64(v.X), float64(v.Y))
		c.Paint()
	}

	surface.WriteToPNG("./test.png")
	surface.Finish()
}

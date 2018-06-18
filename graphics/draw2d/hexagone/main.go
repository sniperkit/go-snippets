package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"code.google.com/p/draw2d/draw2d"

	"github.com/gitchander/heuristic/math/graph2d"
	"github.com/gitchander/heuristic/math/hexm"
)

func saveToPngFile(filePath string, m image.Image) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, m)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Wrote %s OK.\n", filePath)
}

func intMin(x, y int) int {

	if x < y {
		return x
	}

	return y
}

func main() {

	ps := hexm.HexPolygon()

	width := 256
	height := 256

	i := image.NewRGBA(image.Rect(0, 0, width, height))

	gc := draw2d.NewGraphicContext(i)

	var cellRadius = float32(intMin(width, height)) * 0.4

	var centerX = float32(width) * 0.5
	var centerY = float32(height) * 0.5

	T := graph2d.NewTransform()
	T.Scale(cellRadius, cellRadius)
	T.Move(centerX, centerY)

	v := graph2d.Vector{float32(ps[0].X), float32(ps[0].Y)}
	v = T.Apply(v)

	x0 := float64(v.X)
	y0 := float64(v.Y)

	gc.MoveTo(x0, y0)
	for i := 1; i < len(ps); i++ {

		v = graph2d.Vector{float32(ps[i].X), float32(ps[i].Y)}
		v = T.Apply(v)

		x := float64(v.X)
		y := float64(v.Y)

		gc.LineTo(x, y)
	}
	gc.LineTo(x0, y0)

	gc.Stroke()

	saveToPngFile("TestPath.png", i)
}

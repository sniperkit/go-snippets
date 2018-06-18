package main

import (
	"bufio"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"os"

	"github.com/gitchander/go-examples/gif/gol/life"
)

func main() {
	renderFile(pars[1])
}

func renderFile(params RenderParams) {

	t := life.NewTorus(params.Size.Dx, params.Size.Dy)

	ps := PointsFromPattern(params.Pattern)
	life.Points(ps).Move(params.Loc)

	for i := 0; i < 4; i++ {
		ps = life.Nexter3{}.Next(ps)
	}

	SetPoints(t, ps)

	fileName := "test.gif"

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	if err := render(w, t, params); err != nil {
		log.Fatal(err)
	}
}

func render(w io.Writer, t *life.Torus, params RenderParams) error {

	const (
		indexBlack = 0
		indexWhite = 1
	)
	palette := []color.Color{
		indexBlack: color.Black,
		indexWhite: color.White,
	}

	var nX, nY = t.Size()

	var (
		scale = params.Scale
		n     = params.FrameCount
	)

	anim := gif.GIF{LoopCount: n}

	for i := 0; i < n; i++ {
		rect := image.Rect(0, 0, nX*scale, nY*scale)
		ip := image.NewPaletted(rect, palette)

		for y := 0; y < nY; y++ {
			for x := 0; x < nX; x++ {

				indexColor := uint8(indexWhite)
				if t.Get(x, y) {
					indexColor = indexBlack
				}

				rect = image.Rect(x*scale, y*scale, (x+1)*scale, (y+1)*scale)

				fillRect(ip, indexColor, rect)
			}
		}

		t.Next()

		anim.Delay = append(anim.Delay, params.Delay)
		anim.Image = append(anim.Image, ip)
	}

	return gif.EncodeAll(w, &anim)
}

func fillRect(ip *image.Paletted, indexColor uint8, r image.Rectangle) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			ip.SetColorIndex(x, y, indexColor)
		}
	}
}

type RenderParams struct {
	Size       Size
	FrameCount int
	Pattern    Pattern
	Loc        life.Point
	Scale      int
	Delay      int
}

type Size struct {
	Dx int
	Dy int
}

func SetPoints(t *life.Torus, ps life.Points) {
	for _, p := range ps {
		t.Set(int(p.X), int(p.Y), true)
	}
}

var pars = []RenderParams{
	RenderParams{
		Size:       Size{Dx: 15, Dy: 15},
		FrameCount: 8,
		Pattern:    patternGalaxy,
		Loc:        life.Point{3, 3},
		Scale:      30,
		Delay:      8,
	},
	RenderParams{
		Size:       Size{Dx: 42, Dy: 28},
		FrameCount: 30,
		Pattern:    patternGliderGun,
		Loc:        life.Point{3, 3},
		Scale:      8,
		Delay:      8,
	},
}

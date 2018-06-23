package filters

import (
	"appengine"
	"code.google.com/p/draw2d/draw2d"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	_ "image/jpeg"
	"math"
	"math/rand"
)

func generateBrushes(minRad, numBrushes int) []int {
	brushes := make([]int, numBrushes)
	for k := range brushes {
		brushes[numBrushes-k-1] = minRad
		minRad += minRad
	}
	return brushes
}

func FilterPainterly(c appengine.Context, m image.Image) image.Image {
	bounds := m.Bounds()
	canvas := image.NewRGBA(bounds)

	// Estos parámetros posteriormente deberán ser... parametrizados:
	brushMinRadius := 3
	numOfBrushes := 3
	brushes := generateBrushes(brushMinRadius, numOfBrushes)

	for _, radius := range brushes {
		c.Infof("Brush %v", radius)
		refImage := imaging.Blur(m, float64(radius))
		paintLayer(canvas, refImage, radius, 100)
	}
	return canvas
}

type MyStroke struct {
	Color  color.Color
	Point  image.Point
	Radius int
}

func paintLayer(cnv *image.RGBA, refImage image.Image, radius int, T float64) image.Image {
	strokes := make([]MyStroke, 0)
	D := ImageDifference(cnv, refImage)

	ys := cnv.Bounds().Max.Y
	xs := cnv.Bounds().Max.X
	for y := 0; y < ys; y++ {
		for x := 0; x < xs; x++ {
			// Calculates the error near (x,y):
			areaError := float64(0)
			maxdif := float64(0)
			maxx := 0
			maxy := 0
			for y2 := IntMax(0, y-radius); y2 < IntMin(ys, y+radius); y2++ {
				for x2 := IntMax(0, x-radius); x2 < IntMin(xs, x+radius); x2++ {
					dif := D[y2][x2]
					areaError += dif
					if dif > maxdif {
						maxdif = dif
						maxx = x2
						maxy = y2
					}
				}
			}
			areaError = areaError / float64(radius*radius)

			if areaError > T {
				strokes = append(strokes, MyStroke{
					Color:  refImage.At(maxx, maxy),
					Point:  image.Point{maxx, maxy},
					Radius: radius,
				})
			}
		}
	}
	paintStrokes(cnv, strokes)
	return cnv
}

func paintStrokes(cnv *image.RGBA, strokes []MyStroke) {
	gc := draw2d.NewGraphicContext(cnv)
	order := rand.Perm(len(strokes))

	for _, v := range order {
		s := strokes[v]
		gc.SetFillColor(s.Color)
		gc.SetLineWidth(0)
		gc.ArcTo(float64(s.Point.X),
			float64(s.Point.Y),
			float64(s.Radius),
			float64(s.Radius),
			0, 2*math.Pi)
		gc.FillStroke()
	}
}

package filters

import (
	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
	colorful "github.com/lucasb-eyer/go-colorful"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"math/rand"
)

func RescaleImage(m image.Image, size int) image.Image {
	bounds := m.Bounds()
	ys := bounds.Max.Y
	xs := bounds.Max.X
	if xs > ys {
		return imaging.Resize(m, IntMin(size, xs), 0, imaging.Lanczos)
	} else {
		return imaging.Resize(m, 0, IntMin(size, ys), imaging.Lanczos)
	}
}

func bidimensionalArray(x, y int) [][]int {
	res := make([][]int, x)
	for i := range res {
		res[i] = make([]int, y)
	}
	return res
}

func distance(A []int, x, y int) float64 {
	dy := A[1] - y
	dx := A[0] - x
	return math.Sqrt(float64(dy*dy + dx*dx))
}

func manhattan(A []int, x, y int) float64 {
	dy := A[1] - y
	dx := A[0] - x
	return math.Abs(float64(dy)) + math.Abs(float64(dx))
}

func colorMean(colors []color.Color) color.Color {
	var r, g, b, a float64
	r, g, b, a = 0, 0, 0, 0
	for _, v := range colors {
		R, G, B, A := v.RGBA()
		r += float64(R)
		g += float64(G)
		b += float64(B)
		a += float64(A)
	}
	c := float64(len(colors))
	return color.NRGBA{uint8(r / c), uint8(g / c), uint8(b / c), uint8(a / c)}
}

type MyColor struct {
	R, G, B, A, C int64
}

func (o *MyColor) Add(c color.Color) {
	r, g, b, a := c.RGBA()
	o.R += int64(r)
	o.G += int64(g)
	o.B += int64(b)
	o.A += int64(a)
	o.C++
}

func (o *MyColor) Average() color.Color {
	if o.C == 0 {
		return color.Black
	}
	return color.NRGBA64{R: uint16(o.R / o.C), G: uint16(o.G / o.C), B: uint16(o.B / o.C), A: uint16(o.A / o.C)}
}

func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IntMin(a, b int) int {
	if a > b {
		return b
	}
	return a
}

var sobelHCoefficients = []float64{
	1.0, 0.0, -1.0,
	2.0, 0.0, -2.0,
	1.0, 0.0, -1.0,
}
var sobelVCoefficients = []float64{
	1.0, 2.0, 1.0,
	0.0, 0.0, 0.0,
	-1.0, -2.0, -1.0,
}

func SobelV(src image.Image) *image.RGBA {
	// Apply the vertical sobel filter to an image.
	g := gift.New(
		gift.Convolution(
			[]float32{
				1.0, 2.0, 1.0,
				0.0, 0.0, 0.0,
				-1.0, -2.0, -1.0,
			},
			false, false, true, 0.0,
		),
	)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	return dst
}

func SobelH(src image.Image) *image.RGBA {
	// Apply the vertical sobel filter to an image.
	g := gift.New(
		gift.Convolution(
			[]float32{
				1.0, 0.0, -1.0,
				2.0, 0.0, -2.0,
				1.0, 0.0, -1.0,
			},
			false, false, true, 0.0,
		),
	)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	return dst
}

func ColorToGray(c color.Color) uint8 {
	return color.GrayModel.Convert(c).(color.Gray).Y
}

func GradientData(src image.Image) ([][]float64, [][]float64) {
	xs, ys := src.Bounds().Max.X, src.Bounds().Max.Y
	mag := make([][]float64, ys)
	ori := make([][]float64, ys)
	sobelH := SobelH(src)
	sobelV := SobelV(src)
	for y := 0; y < ys; y++ {
		magrow := make([]float64, xs)
		orirow := make([]float64, xs)
		for x := 0; x < xs; x++ {
			GX := float64(ColorToGray(sobelH.At(x, y)))
			GY := float64(ColorToGray(sobelV.At(x, y)))
			magrow[x] = math.Sqrt(GX*GX + GY*GY)
			orirow[x] = math.Atan2(GY, GX)
		}
		mag[y] = magrow
		ori[y] = orirow
	}
	return mag, ori
}

func ColorDistance(A, B color.Color) float64 {
	ar, ag, ab, aa := A.RGBA()
	br, bg, bb, ba := B.RGBA()
	dr, dg, db, da := float64(ar-br), float64(ag-bg), float64(ab-bb), float64(aa-ba)
	return math.Sqrt(dr*dr + dg*dg + db*db + da*da)
}

func ImageDifference(A *image.RGBA, B image.Image) [][]float64 {
	ys := A.Bounds().Max.Y
	xs := A.Bounds().Max.X
	res := make([][]float64, ys)
	for y := 0; y < ys; y++ {
		rowdif := make([]float64, xs)
		for x := 0; x < xs; x++ {
			rowdif[x] = ColorDistance(A.At(x, y), B.At(x, y))
		}
		res[y] = rowdif
	}
	return res
}

// Generates a new color based on a given color and the jitter
// for a painting style.
func RandomizeColor(c color.Color, settings *PainterlySettings) color.NRGBA {
	//sty := settings.Style
	//r,g,b,_ := c.RGBA()
	//return color.NRGBA{uint8(r/255), uint8(g/255), uint8(b/255), uint8(sty.Opacity*255)}
	sty := settings.Style
	r, g, b, _ := c.RGBA()
	R := Clamp64(0, rand.NormFloat64()*sty.JitterRed*65535/2+float64(r), 65535)
	G := Clamp64(0, rand.NormFloat64()*sty.JitterGreen*65535/2+float64(g), 65535)
	B := Clamp64(0, rand.NormFloat64()*sty.JitterBlue*65535/2+float64(b), 65535)

	n := colorful.Color{R / 65535, G / 65535, B / 65535}
	h, s, v := n.Hsv()
	H := Overflow64(0, rand.NormFloat64()*sty.JitterHue*45+h, 360)
	S := Clamp64(0, rand.NormFloat64()*sty.JitterSaturation*0.25+s, 1)
	V := Clamp64(0, rand.NormFloat64()*sty.JitterValue*0.25+v, 1)

	n2 := colorful.Hsv(H, S, V)
	r2, g2, b2 := n2.RGB255()
	return color.NRGBA{r2, g2, b2, uint8(sty.Opacity * 255)}
}

func Clamp64(a, b, c float64) float64 {
	return math.Max(a, math.Min(b, c))
}

func Overflow64(a, b, c float64) float64 {
	if b < a {
		return Overflow64(a, b+360, c)
	} else if b > c {
		return Overflow64(a, b-360, c)
	}
	return b
}

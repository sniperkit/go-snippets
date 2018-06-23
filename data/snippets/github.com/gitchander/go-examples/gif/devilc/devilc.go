package devilc

import (
	"errors"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"io"
	"math"

	"github.com/gitchander/cairo"
)

type CalcFunc func(t float64) Point

var errDevilClosed = errors.New("devil is closed")

func renderScene(canvas *cairo.Canvas, width, height int, a, b float64) {

	var (
		x0 = float64(width+1) * 0.5
		y0 = float64(height+1) * 0.5
	)

	canvasFillColor(canvas, color.White)

	canvas.SetLineWidth(1)
	canvasSetColor(canvas, color.RGBA{127, 127, 127, 255})
	canvas.MoveTo(x0, 0)
	canvas.LineTo(x0, float64(height))
	canvas.MoveTo(0, y0)
	canvas.LineTo(float64(width), y0)
	canvas.Stroke()

	canvasSetColor(canvas, color.RGBA{255, 0, 0, 255})
	canvas.SetLineWidth(2)

	m := cairo.NewMatrix()
	m.InitIdendity()
	m.Translate(x0, y0)
	canvas.Save()
	canvas.SetMatrix(m)
	drawDevilCurve(canvas, a, b)
	canvas.Restore()
}

func drawDevilCurve(canvas *cairo.Canvas, a, b float64) {

	cf := getDevilCalcFunc(a, b)

	const (
		piDiv4 = math.Pi / 4
		d      = 0.00001
		count  = 100
	)
	rs := []Range{
		Range{Min: -piDiv4 + d, Max: piDiv4 - d, Count: count},
		Range{Min: piDiv4 + d, Max: 3*piDiv4 - d, Count: count},
		Range{Min: 3*piDiv4 + d, Max: 5*piDiv4 - d, Count: count},
		Range{Min: 5*piDiv4 + d, Max: 7*piDiv4 - d, Count: count},
	}

	drawCurve(canvas, cf, rs)
}

func drawCurve(c *cairo.Canvas, cf CalcFunc, rs []Range) {
	for _, r := range rs {
		if n := r.Count; n > 1 {
			var (
				t  = minf(r.Min, r.Max)
				dt = r.Step()
			)
			p := cf(t)
			c.MoveTo(p.X, p.Y)
			t += dt
			for i := 1; i < n; i++ {
				p = cf(t)
				c.LineTo(p.X, p.Y)
				t += dt
			}
			c.Stroke()
		}
	}
}

type devilCalcer struct {
	opened  bool
	surface *cairo.Surface
	canvas  *cairo.Canvas
	data    []byte
}

func newDevilCalcer(width, height int) (*devilCalcer, error) {

	surface, err := cairo.NewSurface(cairo.FORMAT_ARGB32, width, height)
	if err != nil {
		return nil, err
	}

	canvas, err := cairo.NewCanvas(surface)
	if err != nil {
		return nil, err
	}

	d := &devilCalcer{
		opened:  true,
		surface: surface,
		canvas:  canvas,
		data:    make([]byte, surface.GetDataLength()),
	}

	return d, nil
}

func (d *devilCalcer) Close() error {
	if !d.opened {
		return errDevilClosed
	}
	d.canvas.Destroy()
	d.surface.Destroy()
	d.opened = false
	return nil
}

func (d *devilCalcer) Calc(a, b float64) error {

	if !d.opened {
		return errDevilClosed
	}

	var (
		width  = d.surface.GetWidth()
		height = d.surface.GetHeight()
	)

	renderScene(d.canvas, width, height, a, b)

	return nil
}

func (d *devilCalcer) Draw(ip *image.Paletted, pal color.Palette) error {

	data := d.data

	err := d.surface.GetData(data)
	if err != nil {
		return err
	}

	indexWhite := uint8(pal.Index(color.White))

	var (
		nX     = d.surface.GetWidth()
		nY     = d.surface.GetHeight()
		stride = d.surface.GetStride()
	)

	const valueWhite uint32 = 0xFFFFFFFF

	for y := 0; y < nY; y++ {
		for x := 0; x < nX; x++ {
			var (
				j = x * 4

				b = data[j]
				g = data[j+1]
				r = data[j+2]
				a = data[j+3]

				u = uint32(b) +
					(uint32(g) << 8) +
					(uint32(r) << 16) +
					(uint32(a) << 24)
			)

			if u == valueWhite {
				ip.SetColorIndex(x, y, indexWhite)
			} else {
				c := color.NRGBA{r, g, b, a}
				ip.Set(x, y, pal.Convert(c))
			}
		}
		data = data[stride:]
	}

	return nil
}

func renderGif(w io.Writer) error {

	pal := palette.WebSafe

	var nX, nY = 512, 512

	var n = 80

	anim := gif.GIF{LoopCount: n}

	dc, err := newDevilCalcer(nX, nY)
	if err != nil {
		return err
	}
	defer dc.Close()

	var (
		mn   = mini(nX, nY)
		minV = float64(mn) / 8
		maxV = float64(mn) / 4
	)

	var parA, parB = minV, maxV
	var dPar = (maxV - minV) / float64(n/2)

	for i := 0; i < n; i++ {

		dc.Calc(parA, parB)
		if i < n/2 {
			parA, parB = parA+dPar, parB-dPar
		} else {
			parA, parB = parA-dPar, parB+dPar
		}

		rect := image.Rect(0, 0, nX, nY)
		ip := image.NewPaletted(rect, pal)

		dc.Draw(ip, pal)

		anim.Delay = append(anim.Delay, 10)
		anim.Image = append(anim.Image, ip)
	}

	return gif.EncodeAll(w, &anim)
}

package devilc

import (
	"image/color"

	"github.com/gitchander/cairo"
)

func colorToFRGBA(c color.Color) (R, G, B, A float64) {
	r, g, b, a := c.RGBA()
	const max = 65535
	R = float64(r) / max
	G = float64(g) / max
	B = float64(b) / max
	A = float64(a) / max
	return
}

func canvasFillColor(canvas *cairo.Canvas, c color.Color) {

	surface := canvas.GetTarget()
	if surface == nil {
		return
	}

	canvas.Save()
	canvas.SetSourceRGBA(colorToFRGBA(c))
	canvas.SetOperator(cairo.OPERATOR_SOURCE)
	canvas.Paint()
	canvas.Restore()
}

func canvasSetColor(canvas *cairo.Canvas, c color.Color) {
	canvas.SetSourceRGBA(colorToFRGBA(c))
}

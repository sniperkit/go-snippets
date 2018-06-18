package main

import (
	"image"
	"image/color"
)

func contrast(v uint32) uint32 {
	vf := float64(v) / 255
	if vf < 0.5 {
		vf = vf * vf * 2
	} else {
		vf = 1 - vf
		vf = 1 - vf*vf*2
	}
	return uint32(255 * vf)

	// if v < 128 {
	// 	return v * v * 2
	// }
	// v = 128 - v
	// return 128 - v*v*2
}

// Convert to binary image
func Transform(src image.Image) image.Image {
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	ret := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := src.At(x, y).RGBA()
			r, g, b = r/256, g/256, b/256

			r = contrast(r)
			g = contrast(g)
			b = contrast(b)

			ret.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 0})
		}
	}
	return ret
}

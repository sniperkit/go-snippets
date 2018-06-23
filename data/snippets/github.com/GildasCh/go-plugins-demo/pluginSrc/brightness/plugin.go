package main

import (
	"image"
	"image/color"
)

func add(v uint32) uint32 {
	v += 50
	if v > 255 {
		v = 255
	}
	return v
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

			r = add(r)
			g = add(g)
			b = add(b)

			ret.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 0})
		}
	}
	return ret
}

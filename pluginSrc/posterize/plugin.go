package main

import (
	"image"
	"image/color"
)

// Convert to binary image
func Transform(src image.Image) image.Image {
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	ret := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := src.At(x, y).RGBA()
			r, g, b = r/256, g/256, b/256
			if r < 100 {
				r = 0
			} else {
				r = 255
			}
			if g < 100 {
				g = 0
			} else {
				g = 255
			}
			if b < 100 {
				b = 0
			} else {
				b = 255
			}
			ret.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 0})
		}
	}
	return ret
}

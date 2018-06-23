package main

import "image"

var Priority int = 1

// Convert to grayscale image
func Transform(src image.Image) image.Image {
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	ret := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			ret.Set(x, y, src.At(x, y))
		}
	}
	return ret
}

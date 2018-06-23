package filters

import (
	"appengine"
	"github.com/disintegration/imaging"
	"image"
)

func FilterGrayscale(_ appengine.Context, m image.Image) image.Image {
	res := imaging.Grayscale(m)
	return res
}

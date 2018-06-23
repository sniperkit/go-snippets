// Stole the majority of this from http://blog.golang.org/go-image-package
package main

// Process for adding an image to a wall:
// 1. Image is uploaded to server
// 2. Image is converted to PNG and normalized
// 3. Image is uploaded to AWS...for some reason
// 4. Image is added to Wall
// 5. Wall regenerates main image, uploads to AWS

// Process for zooming in:
// 1. Get target zoom area
// 2. Do some math or something
// 3. Generate magical new image
// 4. Something?

import (
	"errors"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"io"
	"math"

	_ "code.google.com/p/vp8-go/webp"
	_ "image/jpeg"
)

func Normalize(totalPix int, file io.Reader, w io.Writer) (image.Image, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	if bounds.Dx() == 0 || bounds.Dy() == 0 {
		return nil, errors.New("One or more of your dimensions is zero")
	}
	// Ratio
	ratio := bounds.Dx() / bounds.Dy()
	width := uint(math.Floor(math.Sqrt(float64(ratio * totalPix))))
	img = resize.Resize(width, 0, img, resize.Lanczos3)

	err = png.Encode(w, img)
	if err != nil {
		return nil, err
	}

	return img, nil
}

package main

import (
	_ "code.google.com/p/vp8-go/webp"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
)

func convertToPNG(w io.Writer, r io.Reader) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
}

func convertJPEGToPNG(w io.Writer, r io.Reader) error {
	img, _, err := jpeg.Decode(r)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
}

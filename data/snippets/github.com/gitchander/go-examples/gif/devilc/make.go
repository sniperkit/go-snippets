package devilc

import (
	"bytes"
	"io/ioutil"

	"github.com/gitchander/cairo"
)

func MakePNG(fileName string, width, height int, a, b float64) error {

	surface, err := cairo.NewSurface(cairo.FORMAT_ARGB32, width, height)
	if err != nil {
		return err
	}
	defer surface.Destroy()

	canvas, err := cairo.NewCanvas(surface)
	if err != nil {
		return err
	}
	defer canvas.Destroy()

	renderScene(canvas, width, height, a, b)

	return surface.WriteToPNG(fileName)
}

func MakeGif(fileName string) error {

	var buffer bytes.Buffer

	if err := renderGif(&buffer); err != nil {
		return err
	}

	err := ioutil.WriteFile(fileName, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

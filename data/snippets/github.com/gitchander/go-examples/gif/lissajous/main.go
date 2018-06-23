// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"bufio"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"os"
)

func main() {

	file, err := os.Create("lissajous.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	lissajous(w)
}

func genPalette() color.Palette {
	p := make([]color.Color, 256)
	for i := range p {
		p[i] = color.Gray{Y: uint8(i)}
		//p[i] = color.RGBA{0, uint8(i), 0, 0xFF}
	}
	return p
}

func lissajous(w io.Writer) {

	palette := genPalette()

	const (
		cycles = 2
		// number of complete x oscillator revolutions
		res  = 0.001 // angular resolution
		size = 100
		// image canvas covers [-size..+size]
		nframes = 64
		// number of animation frames
		delay = 8
	// delay between frames in 10ms units
	)

	freq := 3.0 / 2.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)

		maxT := cycles * 2 * math.Pi

		for t := 0.0; t < maxT; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			indexColor := uint8(math.Floor(float64(len(palette)-1) * t / maxT))

			img.SetColorIndex(size+int(x*size*0.9+0.5),
				size+int(y*size*0.9+0.5), indexColor)
		}

		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(w, &anim) // NOTE: ignoring encoding errors
}

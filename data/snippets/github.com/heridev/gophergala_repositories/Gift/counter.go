package gift

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"

	"fmt"
	"image"
	"image/color/palette"
	"io/ioutil"
	"log"
)

// ImageCounter is a quick debug image to print out text over an image.
type ImageCounter struct {
	font          *truetype.Font
	c             *freetype.Context
	width, height int
}

// Geo is a nop
func (g *ImageCounter) Geo(lat, long, heading float64) {
}

// Pipe sends our frames to the server.
func (g *ImageCounter) Pipe(images chan giftImage) {
	defer close(images)
	log.Printf("About to send ImageCounter")
	for i := 0; i < 10; i++ {
		img := image.NewPaletted(image.Rect(0, 0, g.width, g.height), palette.Plan9)
		g.c.SetDst(img)
		pt := freetype.Pt(g.width/2-100, g.height/2)
		g.c.DrawString(fmt.Sprintf("Frame: %d", i), pt)

		images <- giftImage{img: img, frameTimeMS: 100, disposalFlags: disposalRestorePrev}
	}
}

// Setup loads our font
func (g *ImageCounter) Setup(width, height int) {
	fontBytes, err := ioutil.ReadFile("TimesNewRoman.ttf")
	if err != nil {
		log.Println(err)
		return
	}
	g.font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	g.width = width
	g.height = height

	fg := image.White
	g.c = freetype.NewContext()
	g.c.SetDPI(72)
	g.c.SetFont(g.font)
	g.c.SetFontSize(48)
	g.c.SetClip(image.Rect(0, 0, width, height))
	g.c.SetSrc(fg)
}

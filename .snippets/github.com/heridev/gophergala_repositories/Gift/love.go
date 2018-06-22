package gift

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"

	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
)

// ImageLove loads your geo location(or other provided position) and
// gets cupid to fire a arrow at it
type ImageLove struct {
	MapKey, StreetViewKey string

	font          *truetype.Font
	c             *freetype.Context
	lat, long     float64
	width, height int
	httpImages    chan giftImage
}

func (g *ImageLove) drawString(img *image.Paletted, x, y int, text string) int {
	g.c.SetDst(img)
	pt := freetype.Pt(x, y)
	ptAfter, _ := g.c.DrawString(text, pt)
	return int((ptAfter.X >> 8) - (pt.X >> 8))
}

// Geo takes our lat/long and starts streaming some gif images, overlays some
// text on top of the map, along with the targeting reticle, and then plays
// the explosion animation over the final position
func (g *ImageLove) Geo(lat, long, heading float64) {
	g.lat = lat
	g.long = long

	bounds := image.Rect(0, 0, g.width, g.height)

	go func() {
		defer close(g.httpImages)

		fullscreen := image.NewPaletted(bounds, palette.Plan9)
		draw.Src.Draw(fullscreen, bounds, image.White, image.Pt(0, 0))

		g.httpImages <- giftImage{img: fullscreen, frameTimeMS: 0, disposalFlags: disposalNone}

		measure(func() {
			overlayGif("love/cupidarrow.gif", bounds, g.httpImages)
		}, "cupid arrow image")

		strings := []string{
			"",
			"YOU",
			"ARE",
			"LOVED",
			"SO",
			"MUCH",
			"XOXOXO",
		}

		measure(func() {
			startSide := rand.Intn(4)
			var startPt = image.Pt(0, 0)
			switch startSide {
			case 0: // Random Y coord, X = 0
				startPt.X = 0
				startPt.Y = rand.Intn(bounds.Dy())
			case 1: // Random Y coord, X = right side
				startPt.X = bounds.Dx()
				startPt.Y = rand.Intn(bounds.Dy())
			case 2: // Random X coord, Y = 0
				startPt.X = rand.Intn(bounds.Dx())
				startPt.Y = 0
			case 3:
				startPt.X = rand.Intn(bounds.Dx())
				startPt.Y = bounds.Dy()
			}
			for i := 1; i < 7; i++ {
				maptype := "roadmap"
				if i > 4 {
					maptype = "satellite"
				}

				url := mapURL(lat, long, g.width, g.height, i*3, maptype, g.MapKey)
				resp, err := http.Get(url)
				if err != nil {
					log.Printf("Error requesting map: %d: %+v\n", i, err)
					continue
				}

				img, err := gif.Decode(resp.Body)
				if err != nil {
					log.Printf("Error decoding map: %+v", err)
					continue
				}
				g.httpImages <- giftImage{img: img.(*image.Paletted), frameTimeMS: 0}

				img = image.NewPaletted(img.Bounds(), palette.Plan9)
				img.(*image.Paletted).Palette[0] = color.RGBA{0, 0, 0, 0}
				center := image.Pt(img.Bounds().Dx()/2, img.Bounds().Dy()/2)
				width := g.drawString(img.(*image.Paletted), 0, 48, strings[i])

				offset := image.Pt(bounds.Dx()/2-width/2, 0)
				g.httpImages <- giftImage{img: img.(*image.Paletted), frameTimeMS: 10, offset: offset}

				delta := center.Sub(startPt)
				dlen := math.Sqrt(float64(delta.X*delta.X + delta.Y*delta.Y))
				dx, dy := float64(delta.X)/dlen, float64(delta.Y)/dlen

				crosshairSteps := 20
				timeStep := float64(4)
				for j := 0; j < crosshairSteps; j++ {

					t := (float64(j+1) / float64(crosshairSteps)) * dlen

					pt := image.Pt(startPt.X+int(dx*t), startPt.Y+int(dy*t))

					ts := int(timeStep)
					if j == crosshairSteps-1 {
						embedFrame("love/heart.gif", img.Bounds(), pt, disposalRestoreBg, ts+20, g.httpImages)
					} else {
						embedFrame("love/heart_small.gif", img.Bounds(), pt, disposalRestoreBg, ts, g.httpImages)

					}
				}
			}
		}, "google maps queries")

		measure(func() {
			overlayGif("love/explosion.gif", bounds, g.httpImages)
		}, "heart overlay")

		//embedFrame("love/heart_complete.gif", bounds, image.Pt(0, 0), disposalNone, 200, g.httpImages)
	}()
}

// Pipe is a simple pipe between our internal channel, and the channel the server provides
func (g *ImageLove) Pipe(images chan giftImage) {
	log.Printf("About to send love map")
	for pm := range g.httpImages {
		images <- pm
	}
	close(images)
}

// Setup initializes our width and height and loads the font
func (g *ImageLove) Setup(width, height int) {
	g.width = width
	g.height = height

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

	fg := image.NewUniform(color.RGBA{255, 0, 0, 255})
	g.c = freetype.NewContext()
	g.c.SetDPI(72)
	g.c.SetFont(g.font)
	g.c.SetFontSize(48)
	g.c.SetClip(image.Rect(0, 0, width, height))
	g.c.SetSrc(fg)

	g.httpImages = make(chan giftImage)
}

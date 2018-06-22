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
	"os"
	"time"
)

// ImageWar loads your geo location(or other provided position) and
// drops a nuke on it
type ImageWar struct {
	MapKey, StreetViewKey string

	font          *truetype.Font
	c             *freetype.Context
	lat, long     float64
	width, height int
	httpImages    chan giftImage
}

func measure(fun func(), desc string) {
	start := time.Now()
	fun()
	end := time.Now()
	log.Printf("%s took: %+v", desc, end.Sub(start))
}

func overlayGif(path string, bounds image.Rectangle, images chan giftImage) error {
	return embedGif(path, bounds, image.Pt(0, 0), disposalRestoreBg, images)
}

func embedGif(path string, bounds image.Rectangle, offset image.Point, disposalFlags uint8, images chan giftImage) error {
	gifFile, err := os.Open(path)
	if err != nil {
		log.Printf("Unable to open gif: %s", path)
		return err
	}
	defer gifFile.Close()
	frames, err := gif.DecodeAll(gifFile)
	if err != nil {
		log.Printf("Unable to decode gif: %s", path)
		return err
	}

	for i, frame := range frames.Image {

		parentWidth := bounds.Dx()
		parentHeight := bounds.Dy()
		childWidth := frame.Bounds().Dx()
		childHeight := frame.Bounds().Dy()

		offsetPt := image.Pt(parentWidth/2-childWidth/2, parentHeight/2-childHeight/2)

		images <- giftImage{img: frame, frameTimeMS: frames.Delay[i], disposalFlags: disposalRestorePrev, offset: offsetPt}
	}
	return nil
}

func embedFrame(path string, bounds image.Rectangle, offset image.Point, disposalFlags uint8, frameTimeMS int, images chan giftImage) error {
	gifFile, err := os.Open(path)
	if err != nil {
		log.Printf("Unable to open gif: %s", path)
		return err
	}
	defer gifFile.Close()
	frame, err := gif.Decode(gifFile)
	if err != nil {
		log.Printf("Unable to decode gif: %s", path)
		return err
	}

	center := image.Pt(frame.Bounds().Dx()/2, frame.Bounds().Dy()/2)
	offsetPt := image.Pt(0, 0)
	offsetPt.X += offset.X - center.X
	offsetPt.Y += offset.Y - center.Y

	images <- giftImage{img: frame.(*image.Paletted), frameTimeMS: frameTimeMS, disposalFlags: disposalRestorePrev, offset: offsetPt}
	return nil
}

func (g *ImageWar) drawString(img *image.Paletted, x, y int, text string) {
	g.c.SetDst(img)
	pt := freetype.Pt(x, y)
	g.c.DrawString(text, pt)
}

// Geo takes our lat/long and starts streaming some gif images, overlays some
// text on top of the map, along with the targeting reticle, and then plays
// the explosion animation over the final position
func (g *ImageWar) Geo(lat, long, heading float64) {
	g.lat = lat
	g.long = long

	bounds := image.Rect(0, 0, g.width, g.height)

	go func() {
		defer close(g.httpImages)

		measure(func() {
			overlayGif("war/nasr.gif", bounds, g.httpImages)
		}, "rocket launch image")

		strings := []string{
			"",
			"LAUNCH DETECTED",
			"",
			"ACQUIRING TARGET",
			"",
			"TARGET ACQUIRED",
			"",
		}

		measure(func() {
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
				g.httpImages <- giftImage{img: img.(*image.Paletted), frameTimeMS: 10}

				img = image.NewPaletted(img.Bounds(), palette.Plan9)
				img.(*image.Paletted).Palette[0] = color.RGBA{0, 0, 0, 0}
				center := image.Pt(img.Bounds().Dx()/2, img.Bounds().Dy()/2)
				g.drawString(img.(*image.Paletted), 10, 48, strings[i])

				g.httpImages <- giftImage{img: img.(*image.Paletted), frameTimeMS: 10}

				startSide := rand.Intn(4)
				var startPt = image.Pt(0, 0)
				switch startSide {
				case 0: // Random Y coord, X = 0
					startPt.X = 0
					startPt.Y = rand.Intn(img.Bounds().Dy())
				case 1: // Random Y coord, X = right side
					startPt.X = img.Bounds().Dx()
					startPt.Y = rand.Intn(img.Bounds().Dy())
				case 2: // Random X coord, Y = 0
					startPt.X = rand.Intn(img.Bounds().Dx())
					startPt.Y = 0
				case 3:
					startPt.X = rand.Intn(img.Bounds().Dx())
					startPt.Y = img.Bounds().Dy()
				}
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
						embedFrame("war/crosshair.gif", img.Bounds(), pt, disposalRestoreBg, ts+100, g.httpImages)
					} else {
						embedFrame("war/crosshair_small.gif", img.Bounds(), pt, disposalRestoreBg, ts, g.httpImages)

					}
				}
			}
		}, "google maps queries")

		measure(func() {
			overlayGif("war/explosion.gif", bounds, g.httpImages)
		}, "war overlay")

		black := image.NewPaletted(bounds, palette.Plan9)
		draw.Src.Draw(black, bounds, image.Black, image.Pt(0, 0))

		g.httpImages <- giftImage{img: black, frameTimeMS: 200, disposalFlags: disposalLeave}
	}()
}

// Pipe is a simple pipe between our internal channel, and the channel the server provides
func (g *ImageWar) Pipe(images chan giftImage) {
	log.Printf("About to send war map")
	for pm := range g.httpImages {
		images <- pm
	}
	close(images)
}

// Setup initializes our width and height and loads the font
func (g *ImageWar) Setup(width, height int) {
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

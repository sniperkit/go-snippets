package gift

import (
	"github.com/oschwald/geoip2-golang"

	"image"
	"log"
)

// ImageSource should be implemented by our image sources.
// Setup is called first to provide the width/height to make your images
// Geo provides geolocation info
// Pipe is called to give a place to start feeding images to.
type ImageSource interface {
	Setup(width, height int)
	Geo(lat, long, heading float64)
	Pipe(images chan giftImage)
}

// Server is the interface between net/http and the image sources.
type Server struct {
	db     *geoip2.Reader
	width  int
	height int
	source ImageSource
}

type giftImage struct {
	img           *image.Paletted
	frameTimeMS   int
	disposalFlags uint8
	offset        image.Point
}

// NewGiftServer returns a server to be handed to http.HandleFunc.
// Our geoip database is loaded at this point
func NewGiftServer(w, h int, source ImageSource) Server {
	gs := Server{width: w, height: h, source: source}

	var err error
	gs.db, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal("Error opening geoip database: %+v", err)
	}

	return gs
}

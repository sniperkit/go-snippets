package gift

import (
	"fmt"
	"image"
	"image/gif"
	"log"
	"net/http"
)

// ImageMap looks up several map zoom levels from google maps
type ImageMap struct {
	MapKey, StreetViewKey string
	lat, long             float64
	width, height         int
	httpImages            chan giftImage
}

var mapRoot = "https://maps.googleapis.com/maps/api/staticmap"

func mapURL(centerX, centerY float64, width, height, zoom int, maptype string, key string) string {
	var url string
	if key != "" {
		url = fmt.Sprintf("%s?center=%f,%f&zoom=%d&size=%dx%d&format=gif&maptype=%s&key=%s", mapRoot, centerX, centerY, zoom, width, height, maptype, key)
	} else {
		url = fmt.Sprintf("%s?center=%f,%f&zoom=%d&size=%dx%d&format=gif&maptype=%s&", mapRoot, centerX, centerY, zoom, width, height, maptype)
	}
	return url
}

// Geo spawns a goroutine to load our map images and feed them to our internal image channel
func (g *ImageMap) Geo(lat, long, heading float64) {
	g.lat = lat
	g.long = long

	go func() {
		defer close(g.httpImages)

		for i := 1; i < 6; i++ {
			url := mapURL(lat, long, g.width, g.height, i*3, "roadmap", g.MapKey)
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
			g.httpImages <- giftImage{img: img.(*image.Paletted), frameTimeMS: 100}
		}
	}()
}

// Pipe sends our internal images back to the server through the channel
func (g *ImageMap) Pipe(images chan giftImage) {
	defer close(images)
	log.Printf("About to send map")
	for pm := range g.httpImages {
		images <- pm
	}
}

// Setup initializes our ImageMap and creates our internal image channel
func (g *ImageMap) Setup(width, height int) {
	g.width = width
	g.height = height
	g.httpImages = make(chan giftImage)
}

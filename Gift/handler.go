package gift

import (
	"image"
	"image/color/palette"
	"image/gif"
	"log"
	"net"
	"net/http"
	"strconv"
)

// Handler deals with our image sources.  It sets up a pipe between the
// image source, the gif encoder, and the http.ResponseWriter
func (gs *Server) Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request start")

	gs.source.Setup(gs.width, gs.height)

	w.Header().Set("Content-Type", "image/gif")
	// Try our best to never cache this image
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	w.Header().Set("Expires", "0")                                         // Proxies

	var latitude, longitude float64
	var err error

	latCookie, latErr := r.Cookie("latitude")
	longCookie, longErr := r.Cookie("longitude")
	headingCookie, headingErr := r.Cookie("heading")

	latString := r.FormValue("latitude")
	longString := r.FormValue("longitude")
	headingString := "null"

	// If we were able to retrive our cookies, override the strings
	if latErr == nil && longErr == nil && headingErr == nil {
		log.Printf("Attempting cookies")
		latString = latCookie.Value
		longString = longCookie.Value
		headingString = headingCookie.Value
	}

	if latString != "" && longString != "" {
		latitude, err = strconv.ParseFloat(latString, 64)
		if err != nil {
			log.Printf("Unable to parse latitude cookie: %+v", err)
			latitude = 0
		}
		longitude, err = strconv.ParseFloat(longString, 64)
		if err != nil {
			log.Printf("Unable to parse longitude cookie: %+v", err)
			longitude = 0
		}
	} else {
		// Fall back to Geo lookup
		log.Printf("Using geo lookup")

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Error: %+v", err)
		}
		if host == "::1" {
			host = "24.244.32.186"
		}
		ip := net.ParseIP(host)

		record, err := gs.db.City(ip)
		if err != nil {
			log.Printf("Error looking up ip: %+v", err)
		}
		latitude = record.Location.Latitude
		longitude = record.Location.Longitude
	}
	heading, err := strconv.ParseFloat(headingString, 64)
	if err != nil {
		heading = 90
	}
	go gs.source.Geo(latitude, longitude, heading)

	g := gif.GIF{}
	g.Image = append(g.Image, image.NewPaletted(image.Rect(0, 0, gs.width, gs.height), palette.Plan9))
	g.Delay = append(g.Delay, 100)
	g.LoopCount = 0

	images := make(chan giftImage)

	go gs.source.Pipe(images)

	err = EncodeAll(w, &g, images)
	if err != nil {
		log.Printf("Err: %+v", err)
	}
	log.Printf("Request complete")

}

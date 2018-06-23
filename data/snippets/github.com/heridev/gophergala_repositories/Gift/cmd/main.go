package main

import (
	"github.com/gophergala/Gift"

	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

func storeGeoInCookies(w http.ResponseWriter, r *http.Request) {
	var body string
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		body = string(b)
	}

	u, err := url.Parse("?" + body)
	if err != nil {
		log.Printf("Error parsing setgeo body: %+v", err)
	}
	values := u.Query()

	cookie := http.Cookie{Name: "latitude", Value: values.Get("latitude"), Path: "/"}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "longitude", Value: values.Get("longitude"), Path: "/"}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "heading", Value: values.Get("heading"), Path: "/"}
	http.SetCookie(w, &cookie)

	log.Printf("Setting geo cookies to: [%s, %s] @ %s", values.Get("latitude"), values.Get("longitude"), values.Get("heading"))

}

func main() {
	staticMapKey := os.Getenv("STATIC_MAP_KEY")
	streetViewKey := os.Getenv("STREET_VIEW_KEY")

	rand.Seed(time.Now().UTC().UnixNano())
	log.Printf("Starting GIFT server")

	http.HandleFunc("/counter.gif", func(w http.ResponseWriter, r *http.Request) {
		counterGiftServer := gift.NewGiftServer(640, 480, &gift.ImageCounter{})
		counterGiftServer.Handler(w, r)
	})
	http.HandleFunc("/map.gif", func(w http.ResponseWriter, r *http.Request) {
		mapGiftServer := gift.NewGiftServer(640, 480, &gift.ImageMap{MapKey: staticMapKey, StreetViewKey: streetViewKey})
		mapGiftServer.Handler(w, r)
	})
	http.HandleFunc("/war.gif", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("WAR")
		warGiftServer := gift.NewGiftServer(480, 360, &gift.ImageWar{MapKey: staticMapKey, StreetViewKey: streetViewKey})
		warGiftServer.Handler(w, r)
	})
	http.HandleFunc("/love.gif", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("LOVE")
		loveGiftServer := gift.NewGiftServer(480, 360, &gift.ImageLove{MapKey: staticMapKey, StreetViewKey: streetViewKey})
		loveGiftServer.Handler(w, r)
	})
	http.HandleFunc("/setgeo", storeGeoInCookies)
	http.Handle("/", http.FileServer(http.Dir("static")))

	http.ListenAndServe(":8080", nil)
}

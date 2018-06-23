package hello

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"time"

	"appengine"
	"appengine/urlfetch"
)

func init() {
	m := martini.Classic()
	// Middle ware stuff
	m.Use(render.Renderer(render.Options{
		Directory: "public",
	}))
	m.Use(martini.Logger())
	m.Use(martini.Static("public"))

	// Request handlers
	m.Get("/", func(r render.Render, req *http.Request) {
		r.HTML(200, "index", nil)
	})

	m.Get("/api", func(r render.Render, req *http.Request) {
		qs := req.URL.Query()
		log.Println("Hitting API endpoint")
		c := appengine.NewContext(req)

		gif := queryNext(qs, c)
		log.Println(gif)

		content := gif.Content
		r.JSON(200, map[string]interface{}{
			"id":      gif.GifId,
			"title":   gif.GifTitle,
			"game":    gif.GameTitle,
			"upvotes": content.Upvotes,
			"reddit":  content.Comments,
		})
	})

	m.Get("/scrape", func(res http.ResponseWriter, req *http.Request) {
		c := appengine.NewContext(req)
		client := urlfetch.Client(c)
		var g chan Gif = make(chan Gif)
		scrapeSubreddit("smashbros", client, g)

		isDone := false
		current := time.Now()

		for !isDone {
			select {
			case gif := <-g:
				storeGif(gif, c)
			case <-time.After(time.Second * 60):
				isDone = true
			}
		}

		after := time.Now()

		log.Println("DONE SCRAPING")
		log.Printf("Scraping took %v to run.\n", after.Sub(current))
		res.WriteHeader(200)
	})

	http.Handle("/", m)
}

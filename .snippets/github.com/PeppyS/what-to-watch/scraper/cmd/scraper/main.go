package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PeppyS/what-to-watch/scraper/api"
	"github.com/PeppyS/what-to-watch/scraper/imdb"
	"github.com/PeppyS/what-to-watch/scraper/rottentomatoes"
	"github.com/gocolly/colly"
)

var (
	APIURL = os.Getenv("API_URL")
)

func main() {
	imdbScraper := imdb.NewScraper(
		colly.NewCollector(
			colly.Async(true),
		),
	)
	rottenTomatoesScraper := rottentomatoes.NewScraper(&http.Client{})

	imdbMovies, err := imdbScraper.Scrape()
	if err != nil {
		log.Fatalln("Problem scraping for IMDB movies:", err)
	}

	rottenTomatoesMovies, err := rottenTomatoesScraper.Scrape()
	if err != nil {
		log.Fatalln("Problem scraping for rotten tomatoes movies:", err)
	}

	log.Println("Successfully scraped movies from IMDB & Rotten Tomatoes")

	apiClient := api.NewClient(http.DefaultClient, APIURL)
	err = apiClient.NormalizeAndSend(imdbMovies, rottenTomatoesMovies)
	if err != nil {
		log.Fatalln("Problem normalzing and sending movies to API:", err)
	}

	log.Println("Successfully finished scraping and sending movies to API")
}

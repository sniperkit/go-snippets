package rottentomatoes

import (
	"encoding/json"
	"net/http"
)

const (
	moviesInTheatersAPIEndpoint = "https://www.rottentomatoes.com/api/private/v2.0/browse?maxTomato=100&maxPopcorn=100&services=amazon%3Bhbo_go%3Bitunes%3Bnetflix_iw%3Bvudu%3Bamazon_prime%3Bfandango_now&certified&sortBy=popularity&type=in-theaters&page=1"
)

type Movie struct {
	Title              string
	TomatoScore        int
	PopcornScore       int
	TheaterReleaseDate string
	MpaaRating         string
	Synopsis           string
	SynopsisType       string
	Runtime            string
	Posters            struct {
		Primary string
	}
}

type Scraper struct {
	httpClient *http.Client
}

func NewScraper(c *http.Client) *Scraper {
	return &Scraper{c}
}

func (s *Scraper) Scrape() ([]Movie, error) {
	var movies []Movie

	resp, err := s.httpClient.Get(moviesInTheatersAPIEndpoint)
	if err != nil {
		return movies, err
	}

	var moviesAPIResponse struct {
		Results []Movie
	}

	err = json.NewDecoder(resp.Body).Decode(&moviesAPIResponse)
	if err != nil {
		return movies, err
	}

	return moviesAPIResponse.Results, nil
}

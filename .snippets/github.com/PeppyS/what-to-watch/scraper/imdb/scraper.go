package imdb

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const (
	moviesInTheatersPage = "https://www.imdb.com/showtimes/location"

	linkSelector = "div.lister-item-image > a[href]:nth-child(1)"

	moviePageSelector   = "div.title-overview"
	titleSelector       = "div.title_wrapper h1"
	yearSelector        = "span#titleYear"
	ratingSelector      = "div.title_wrapper div meta[itemprop=contentRating]"
	movieRatingSelector = "div.imdbRating span[itemprop=ratingValue]"
	genreSelector       = "div.title_wrapper a span"
	imageSelector       = "div.poster > a > img"
)

type Movie struct {
	Title       string
	Rating      string
	MovieRating float64
	Genre       string
	Image 		string
}

type Scraper struct {
	*colly.Collector
}

func NewScraper(c *colly.Collector) *Scraper {
	return &Scraper{c}
}

func (s *Scraper) Scrape() ([]Movie, error) {
	var movies []Movie

	// Crawl all movies in theaters
	s.OnHTML(linkSelector, func(e *colly.HTMLElement) {
		e.Request.Visit(strings.Replace(e.Attr("href"), "/showtimes", "", 1))
	})

	// Scrape movie info
	s.OnHTML(moviePageSelector, func(e *colly.HTMLElement) {
		title := strings.TrimSpace(
			// Remove year from title
			strings.Replace(e.ChildText(titleSelector), e.ChildText(yearSelector), "", 1),
		)
		rating := strings.TrimSpace(e.ChildAttr(ratingSelector, "content"))
		movieRating, _ := strconv.ParseFloat(strings.TrimSpace(e.ChildText(movieRatingSelector)), 64)
		genre := strings.TrimSpace(e.ChildText(genreSelector))
		image := e.ChildAttr(imageSelector, "src")

		movies = append(movies, Movie{
			title,
			rating,
			movieRating,
			genre,
			image,
		})
	})

	err := s.Visit(moviesInTheatersPage)
	if err != nil {
		return nil, err
	}

	s.Wait()

	return movies, nil
}

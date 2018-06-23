// Abstracts a lot of the scraping stuff
// for ideally any subreddit
package hello

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	GFYCAT_SUFFIX = map[string]string{
		"q": "site%3Agfycat.com",
	}

	YOUTUBE_SUFFIX = map[string]string{
		"q": "site%3Ayoutube.com",
	}

	QUERY_SEARCH_PARAMS = map[string]string{
		"restrict_sr": "on",  // restrict subreddit
		"t":           "all", // all time
	}

	SORT_OPTIONS = [...]string{
		"relevance",
		"new",
		"hot",
		"top",
		"comments",
	}
	re = regexp.MustCompile(`^https?:\/\/[a-z\:0-9.]+\/`)
)

func prepareUrl(base string, extraKey string, extraVal string) string {
	var queryParams []string

	for k, v := range QUERY_SEARCH_PARAMS {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", k, v))
	}
	for k, v := range GFYCAT_SUFFIX {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", k, v))
	}
	queryParams = append(queryParams, fmt.Sprintf("%s=%s", extraKey, extraVal))

	return fmt.Sprintf("%s/search?%s", base, strings.Join(queryParams, "&"))
}

// Called during the init, called to fetch all the data
// and store in an abstracted data storage
func scrapeSubreddit(name string, client *http.Client, g chan Gif) {
	BASE_URL := fmt.Sprintf("https://reddit.com/r/%s", name)
	log.Println("Scraping the ROOT!!")

	for _, v := range SORT_OPTIONS {
		go scrapeRoot(prepareUrl(BASE_URL, "sort", v), client, g)
		time.Sleep(time.Second * time.Duration(15))
	}
}

// Given the first page of the page, scrape until
// there is no more next button
func scrapeRoot(url string, client *http.Client, g chan Gif) {
	log.Println("Scraping root at : ", url)
	depth := 100
	for nextUrl := url; nextUrl != "" && depth > 0; depth -= 1 {
		log.Println("Scraping next URL", depth, nextUrl)
		time.Sleep(time.Second * time.Duration(4))
		nextUrl = scrapePage(nextUrl, client, g)
	}
	log.Println("The querying is done for ", url)
}

func scrapePage(url string, client *http.Client, g chan Gif) string {
	resp, err := client.Get(url)
	check(err)

	doc, err := goquery.NewDocumentFromResponse(resp)
	check(err)

	doc.Find(".linkflair").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Find("a.thumbnail").Attr("href")
		if !exists {
			return
		}

		// Gets all of the data
		votes, err := strconv.Atoi(s.Find("div.score.unvoted").Text())
		if err != nil || votes < 20 {
			return
		}

		titles := s.Find("p.title").Children().First()
		comments, exists := s.Find(".comments").Attr("href")
		if !exists {
			return
		}

		// We're targetting smash bros for now
		gameTitle := titles.Text()
		gifTitle := titles.Next().Text()

		link = re.ReplaceAllString(link, "")
		gifId := strings.Split(link, "?")[0]

		// Yield the Gif that was just scraped
		g <- Gif{
			Content{
				Comments:  comments,
				Upvotes:   votes,
				Subreddit: "smashbros",
			},
			gameTitle,
			gifTitle,
			gifId,
		}

		// log.Println(gifId)
	})

	nextUrl := ""
	doc.Find("span.nextprev").Children().Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "next") {
			nextUrl, _ = s.Attr("href")
		} else {
			log.Println(s.Text())
		}
	})

	if nextUrl == "" {
		log.Println(doc.Find("body").Text())
	}
	log.Println("Next URL: ", nextUrl)
	return nextUrl
}

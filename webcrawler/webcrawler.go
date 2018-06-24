package webcrawler

import (
	"fmt"
	"github.com/rsmohamad/comp4321/database"
	"github.com/rsmohamad/comp4321/models"
	"net/url"
	"sync"

	"github.com/temoto/robotstxt"
)

// HashMap for storing the robots data for each hostname
var robotMap sync.Map

func fetchRobotsTxt(urlObject *url.URL) (rv *robotstxt.RobotsData) {
	// If there's no response or timeout exceeded, give the default robots checker
	// which will allow all paths within that website
	rv = &robotstxt.RobotsData{}
	robotUrl := fmt.Sprintf("%s://%s/robots.txt", urlObject.Scheme, urlObject.Host)
	res, _ := robotClient.Get(robotUrl)

	// Try parsing the robots.txt
	if res != nil && res.StatusCode == 200 {
		robots, err := robotstxt.FromResponse(res)
		if err == nil {
			rv = robots
			fmt.Println("Fetched " + robotUrl)
		}
		res.Body.Close()
	}
	return
}

// Checks if the url is crawlable.
// Will fetch robots.txt if not previously fetched.
// Thread safe.
func isAllowedToCrawl(link string) bool {
	// Use url object to process URLs
	urlObject, _ := url.Parse(link)
	var robots *robotstxt.RobotsData

	// Fetch a website's robots.txt if it's not already fetched
	res, found := robotMap.Load(urlObject.Host)
	if !found {
		// Ensure that the same robots.txt is not fetched twice
		robotMap.Store(urlObject.Host, nil)
		robots = fetchRobotsTxt(urlObject)
		robotMap.Store(urlObject.Host, robots)
	} else {
		// Wait for other thread to get the robots.txt
		for res == nil {
			res, _ = robotMap.Load(urlObject.Host)
		}
		robots = res.(*robotstxt.RobotsData)
	}

	// Check with robots.txt
	if !robots.TestAgent(urlObject.Path, "Agent") {
		fmt.Println(urlObject.String() + " is not allowed; skipped!")
		return false
	}

	return true
}

// Concurrent routine for fetching a page.
// Feeds the page to results channel if fetch is successful.
// Feeds nil if fetch is unsuccessful.
func concurrentFetch(url string, results *chan *models.Document) {
	if !isAllowedToCrawl(url) {
		*results <- nil
		return
	}
	<-throttle
	page := Fetch(url)
	if page != nil && len(page.Words) > 0 && page.Title != "" {
		*results <- page
	} else {
		*results <- nil
	}
}

func Crawl(uri string, num int, index *database.Indexer, restrictHost, aggressive bool) (pages []*models.Document) {
	var activeCounter int
	var updateWg sync.WaitGroup
	visited := make(map[string]bool)
	results := make(chan *models.Document)
	queue := make([]string, 0)

	initClients(aggressive)

	// Sanitize url
	u, _ := url.Parse(uri)
	newUrl := "http://" + u.Host + u.Path
	fmt.Println(newUrl)
	queue = append(queue, newUrl)
	visited[newUrl] = true

	for len(pages) < num {
		// Create goroutines as needed
		needed := num - len(pages) - activeCounter
		for ; len(queue) > 0 && needed > 0; needed-- {
			activeCounter++
			go concurrentFetch(queue[0], &results)
			queue = queue[1:]
		}

		// End prematurely if no links are available
		if activeCounter <= 0 {
			if len(queue) == 0 {
				break
			} else {
				continue
			}
		}

		// Retrieve one page from results channel
		page := <-results
		activeCounter--
		if page == nil {
			continue
		}

		visited[page.Uri] = true
		pages = append(pages, page)
		fmt.Printf("Fetched page #%d out of %d : %s\n", len(pages), num, page.Uri)
		updateWg.Add(1)
		go func(i int, doc *models.Document) {
			index.UpdateOrAddPage(doc)
			//fmt.Printf("Indexed page #%d out of %d : %s\n", i, num, page.Uri)
			updateWg.Done()
		}(len(pages), page)

		// Put unvisited links into queue
		for _, link := range page.Links {
			// Crawl only the starting domain
			if restrictHost {
				urlParse, _ := url.Parse(link)
				if urlParse.Host != u.Host {
					continue
				}
			}

			// skip if the link is already visited or indexed
			if visited[link] || index.ContainsUrl(link) {
				continue
			}

			queue = append(queue, link)
			visited[link] = true
		}
	}

	updateWg.Wait()
	index.FlushInverted()
	return pages
}

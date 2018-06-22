package crawler

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/GopherGala/i_love_indexes/conn_throttler"
	"github.com/GopherGala/i_love_indexes/elasticsearch"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/errgo.v1"
)

var (
	nginxServerRegexp      = regexp.MustCompile(`^nginx.*$`)
	apacheServerRegexp     = regexp.MustCompile(`^Apache.*$`)
	lighthttpdServerRegexp = regexp.MustCompile(`^lighttpd.*$`)
)

type Crawler interface {
	Crawl() error
}

func NewCrawler(indexOf *elasticsearch.IndexOf, path string) (Crawler, error) {
	sem := conn_throttler.Acquire(indexOf.Host)
	log.Println("Get Index Of:", indexOf.URL()+path)
	res, err := http.Get(indexOf.URL() + path)
	if err != nil {
		sem.Release()
		return nil, errgo.Mask(err)
	}
	if res.StatusCode != 200 {
		b, err := ioutil.ReadAll(res.Body)
		log.Println(string(b), err)
		res.Body.Close()
		sem.Release()
		return nil, errgo.Newf("invalid status code: %v", res.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		res.Body.Close()
		sem.Release()
		return nil, errgo.Mask(err)
	}

	res.Body.Close()
	sem.Release()

	server := res.Header.Get("Server")
	baseCrawler := BaseCrawler{
		relativePath: path,
		itemsToIndex: make(chan *elasticsearch.IndexItem, 10),
		IndexOf:      indexOf,
		Doc:          doc,
	}

	log.Println("Start crawler of:", indexOf.URL()+path)
	go baseCrawler.Start()

	// check text if apache is reverse-proxied
	if apacheServerRegexp.MatchString(server) || strings.Contains(doc.Text(), "Apache/2.") {
		return &ApacheCrawler{baseCrawler}, nil
	} else if lighthttpdServerRegexp.MatchString(server) {
		return &LighttpdCrawler{baseCrawler}, nil
	} else if nginxServerRegexp.MatchString(server) {
		return &NginxCrawler{baseCrawler}, nil
	} else {
		return nil, errgo.Newf("Unknown 'Server' header: %v", server)
	}
}

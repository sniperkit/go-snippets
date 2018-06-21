package videos

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/proxy"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Video struct {
	gorm.Model
	Url        string
	ImageSrc   string
	Name       string
	Alias      string
	Desc       string
	VideoSrc   string
	ThunderSrc string
}

var v *Video
var cRounds int = 0
var dRounds int = 0

const dbDest = "80s-20171226.db"
const outFile = "proxys.txt"

var proxyLisManual = []string{"http://1.28.144.19:80", "http://42.245.252.35:80",
	"http://42.245.252.36:80", "http://111.6.187.52:80", "http://1.28.144.12:80",
	"http://124.232.163.4:3128", "http://85.15.176.34:3128"}

func main() {
	dbInit()
	db, err := gorm.Open("sqlite3", dbDest)
	checkErr(err)
	defer db.Close()
	crawl(db)
}

func dbInit() {
	db, err := gorm.Open("sqlite3", dbDest)
	checkErr(err)
	defer db.Close()
	db.AutoMigrate(&Video{})
}

func readProxy() []string {
	var proxyLis []string
	f, err := os.Open(outFile)
	checkErr(err)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		proxyLis = append(proxyLis, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return proxyLis
}

func crawl(db *gorm.DB) {
	c := colly.NewCollector()
	// c.CacheDir = "./80s2.cache"
	host := "http://www.80s.tw"
	c.SetDebugger(&debug.LogDebugger{})
	c.AllowedDomains = []string{"www.80s.tw"}
	c.MaxDepth = 1
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 20,
		Delay:       5 * time.Second,
	})

	c.URLFilters = []*regexp.Regexp{
		regexp.MustCompile("http://www.80s\\.tw/movie/.*"),
	}

	proxys := readProxy()
	rp, err := proxy.RoundRobinProxySwitcher(proxys...)
	checkErr(err)
	c.SetProxyFunc(rp)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36"

	detailCollector := c.Clone()
	detailCollector.OnHTML("body div#body", func(e *colly.HTMLElement) {
		v = &Video{}
		url := e.Request.URL
		imgSrc := e.ChildAttr("img", "src")
		name := e.ChildText("h1")
		// desc := e.ChildText("div.info span:first-child")
		// esc := e.DOM.Find("div.info span").Get(1)

		desc := childOfText(e.DOM, 1, "div.info span")
		alias := childOfText(e.DOM, 2, "div.info span")
		aliasSlice := strings.Split(alias, " ")[1:]
		alias = strings.Join(aliasSlice, "")
		// content := e.ChildText("div#movie_content")
		videoSrc := e.ChildAttr("input[name='list[]']", "value")
		thunderLink := e.ChildAttr("span.dlname span a", "href")
		v.Url = fmt.Sprint(url)
		v.ImageSrc = imgSrc
		v.Name = name
		v.Desc = desc
		v.Alias = alias
		v.VideoSrc = videoSrc
		v.ThunderSrc = thunderLink
		db.Create(v)

	})

	c.OnHTML("ul.me1 li", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		detailLink := host + link

		detailCollector.Visit(detailLink)
	})

	c.OnHTML("div.pager", func(e *colly.HTMLElement) {
		childs := e.DOM.Children()
		length := childs.Length()
		link := nextPage(e.DOM, length-3)
		if len(link) > 0 {
			nextPage := host + link
			fmt.Println("found next page link: ", nextPage)
			c.Visit(nextPage)
		} else {
			fmt.Println("have no next")
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Status code: ", r.StatusCode)
		fmt.Println("Request URL: ", r.Request.URL, "faild with response: ", string(r.Body[:]), "\nERROR: ", err)
		TryAgain(r, 3, "page")

	})

	detailCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Status code: ", r.StatusCode)
		fmt.Println("Request URL: ", r.Request.URL, "faild with response: ", string(r.Body[:]), "\nERROR: ", err)
		// TryAgain(r, 3, "detail")
	})

	c.Visit("http://www.80s.tw/movie/list")
	c.Wait()
	detailCollector.Wait()
}

func TryAgain(res *colly.Response, round int, mode string) {
	var p *int
	if mode == "page" {
		p = &cRounds
	} else {
		p = &dRounds
	}

	for *p < round {
		*p = *p + 1
		fmt.Println("URL: ", res.Request.URL, "\nTry: ", *p, " times")
		sleepTime := time.Duration((round*2)+3) * time.Second
		time.Sleep(sleepTime)
		err := res.Request.Retry()
		if err != nil {
			if *p == 3 {
				fmt.Println("URL:", res.Request.URL, " scrape retry failed! 3 times")
				*p = 0
				break
			}
		} else {
			fmt.Println("URL:", res.Request.URL, " scrape retry successfully! at rounds: ", *p)
			*p = 0
			break
		}
	}
}

func SaveImage(name, url string) (err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	checkErr(err)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36")
	resp, err := client.Do(req)
	checkErr(err)

	defer resp.Body.Close()
	saveName := "imgs/" + name + ".jpg"
	imgFile, err := os.Create(saveName)
	checkErr(err)

	_, err = io.Copy(imgFile, resp.Body)
	checkErr(err)
	imgFile.Close()
	fmt.Printf("image: %s saved\n", name)
	return
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func childAttr(ele *goquery.Selection, goquerySelector, attrName string) string {
	if attr, ok := ele.Find(goquerySelector).Attr(attrName); ok {
		return strings.TrimSpace(attr)
	}
	return ""
}

func childText(ele *goquery.Selection, goquerySelector string) string {
	return strings.TrimSpace(ele.Find(goquerySelector).Text())
}

func childOfText(ele *goquery.Selection, index int, goquerySelector string) string {
	ret := ""
	ele.Find(goquerySelector).Each(func(i int, s *goquery.Selection) {
		if i == index {
			ret = strings.TrimSpace(s.Text())
		}
	})
	return ret
}

func nextPage(ele *goquery.Selection, index int) string {
	var ok bool
	ret := ""
	length := index + 3
	ele.Find("a").Each(func(i int, s *goquery.Selection) {
		if i == index {
			if ret, ok = s.Attr("href"); ok {
				ret = strings.TrimSpace(ret)
			}
		}
		if length == index {
			if ret, ok = s.Attr("href"); ok {
				ret = ""
			}
		}
	})
	return ret
}

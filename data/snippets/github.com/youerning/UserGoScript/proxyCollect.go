package videos

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/proxy"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
	"net/url"
	// "strconv"
	"bufio"
	// "io/ioutil"
	"os"
	"strings"
	"time"
)

type Proxy struct {
	gorm.Model
	IP         string
	Port       string
	Location   string `gorm:"index"`
	LastUpdate string
	Mode       string
}

const dbDest = "proxy.db"
const checkTimeout = 2 * time.Second
const outFile = "proxys.txt"
const checkDomain = "http://www.80s.tw"

var ss = []string{"socks5://127.0.0.1:1080"}
var proxyLisManual = []string{"socks5://127.0.0.1:1080"}
var proxyLis []string

// var proxyClient *http.Client

func main() {
	if len(os.Args) == 1 {
		fmt.Println("请指定操作类型: crawl, check, output")
		os.Exit(1)
	}
	option := os.Args[1]
	switch option {
	case "crawl":
		crawl()
	case "check":
		checkProxy()
		// checkProxy(proxyLisManual)
	case "read":
		readProxy()
	default:
		fmt.Println("请指定操作类型: crawl, check, output")
		os.Exit(1)
	}
}

func dbInit() {
	db, err := gorm.Open("sqlite3", dbDest)
	checkErr(err)
	defer db.Close()
	db.AutoMigrate(&Proxy{})
}

func readProxy() {
	// checkProxy()
	f, err := os.Open(outFile)
	checkErr(err)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func checkProxy(proxys ...[]string) {
	dbInit()
	db, err := gorm.Open("sqlite3", dbDest)
	checkErr(err)
	defer db.Close()
	var tmpRet []map[string]bool
	if len(proxys) == 0 {
		var rows []Proxy
		db.Find(&rows)
		urlLength := len(rows)
		proxyChan := make(chan string, urlLength)
		result := make(chan map[string]bool, 100)
		for _, r := range rows {
			u := fmt.Sprintf("http://%s:%s", r.IP, r.Port)
			proxyChan <- u
		}
		close(proxyChan)

		for w := 0; w < 30; w++ {
			go proxyVerify(w, proxyChan, result)
		}
		fmt.Println("======>", urlLength)
		for k := 0; k < urlLength; k++ {
			// fmt.Println("ret: ", k)
			tmpRet = append(tmpRet, <-result)
			// fmt.Println(tmpRet)
		}

		success := 0
		f, err := os.Create(outFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		for _, pxy := range tmpRet {
			for k, v := range pxy {
				if v {
					success += 1
					// fmt.Println(k
					_, err := f.WriteString(k + "\n")
					checkErr(err)
				}
			}
		}
		f.Sync()

		fmt.Println("Proxy total:", urlLength)
		fmt.Printf("Find vaild proxy: %v entries", success)

	} else {
		for _, proxyUrl := range proxys[0] {
			fmt.Println("URL: ", proxyUrl)
			Url, err := url.Parse(proxyUrl)
			if err != nil {
				fmt.Println("==>url Parse error <==")
				fmt.Println("URL:", proxyUrl, err)
			}
			proxyClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(Url)},
				Timeout: checkTimeout}
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("==> Get panic: ", err)
				}
			}()
			if res, err := proxyClient.Head(checkDomain); err != nil {
				fmt.Println("test faild!!!\nError:", err)
			} else {
				// fmt.Println("URL: ", proxy, "==>", res.Status)
				if res.StatusCode == 200 {
					fmt.Println("test success!!! ")
				} else {
					fmt.Println("test faild\nStatus: ", res.Status)
				}

			}
		}
	}

}

func proxyVerify(w int, proxyChan <-chan string, result chan<- map[string]bool) {
	for proxy := range proxyChan {
		fmt.Println("woker:", w, "verfify url:", proxy)
		Url, err := url.Parse(proxy)
		if err != nil {
			fmt.Println("==>url Parse error <==")
			fmt.Println("URL:", proxy, err)
		}
		proxyClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(Url)},
			Timeout: checkTimeout}
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("==> Get panic: ", err)
			}
		}()
		if res, err := proxyClient.Head(checkDomain); err != nil {
			result <- map[string]bool{proxy: false}
		} else {
			// fmt.Println("URL: ", proxy, "==>", res.Status)
			if res.StatusCode == 200 {
				result <- map[string]bool{proxy: true}
			} else {
				result <- map[string]bool{proxy: false}
			}

		}
	}

}

func crawl() {
	dbInit()
	db, err := gorm.Open("sqlite3", dbDest)
	checkErr(err)
	defer db.Close()
	c := colly.NewCollector()
	// c.CacheDir = "./kuaidaili.cache"
	// Url := "http://www.kuaidaili.com"
	c.SetDebugger(&debug.LogDebugger{})
	c.AllowedDomains = []string{"free-proxy-list.net"}
	c.MaxDepth = 1
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
		Delay:       3 * time.Second,
	})

	rp, err := proxy.RoundRobinProxySwitcher(ss...)
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36"

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		p := &Proxy{}
		// timeFmt := "2006-01-02 15:04:05"
		dom := e.DOM
		ip := childOfText(dom, 0, "td")
		port := childOfText(dom, 1, "td")
		localtion := childOfText(dom, 3, "td")
		mode := childOfText(dom, 4, "td")
		lastUpdate := childOfText(dom, 7, "td")

		p.IP = ip
		p.Port = port
		p.Mode = mode
		p.Location = localtion
		p.LastUpdate = lastUpdate
		db.Create(p)
		fmt.Println(ip, port, mode, localtion, lastUpdate)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL: ", r.Request.URL, "faild with response: ", r.Body, "\nERROR: ", err)

	})

	// initURl := fmt.Sprintf("http://www.kuaidaili.com/free/inha/%s/", strconv.Itoa(page))
	initUrl := "https://free-proxy-list.net/"
	c.Visit(initUrl)
	c.Wait()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
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

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	// URL "net/url"
	"bufio"
	"os"
	"regexp"
	// "strings"
	"errors"
	"runtime"
	"time"
)

func main() {
	var urlChan = make(chan string, 100)
	urlChan <- startUrl
	contChan := downWorker(urlChan)
	// contChan := downWorker(urlChan)
	retChan := parseWorker(urlChan, contChan)

	wf, _ := os.Create("tt.txt")
	defer wf.Close()
	w := bufio.NewWriter(wf)

	for i := 0; i < 100; i++ {
		p("wating for ret....")
		p("urlChan===> ", len(urlChan))
		p("contChan==> ", len(contChan))
		ret := <-retChan
		fmt.Println("magnet recived.....")
		if i == 99 {
			done = true
		}

		p("coroutine>>>>>>", runtime.NumGoroutine())
		p("url chanel>>>>> ", len(urlChan))
		line := fmt.Sprintf("Magnet:%s\n", ret[:len(ret)])
		p(line)
		w.WriteString(line)
		w.Flush()
	}
}

var p = fmt.Println
var startUrl = "http://www.btmeet.org/search/%E6%80%A7%E6%84%9F.html"
var host = "http://www.btmeet.org"
var visted = map[string]bool{}

// var uRex, _ = regexp.Compile(startUrl + "\\w{1,4}/" + ".*?\\.html")
// var mRex, _ = regexp.Compile(`(magnet\:\?.+?)&`)
var uRex, _ = regexp.Compile("href=\"/(wiki){0,1}(search){0,1}/.+?\\.html")
var mRex, _ = regexp.Compile(`(magnet\:\?.+?)"`)
var done bool = false

func parse(cont string) (mLis []string, uLis []string, err error) {
	mLis = mRex.FindAllString(cont, -1)
	uLis = uRex.FindAllString(cont, -1)
	// p(uLis)
	if len(mLis) == 0 && len(uLis) == 0 {
		err = errors.New("cont don't contain any magnet or url")
	}

	return
}

func httpGet(url string) (content string, status int) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", getUserAgent())
	tr := &http.Transport{IdleConnTimeout: 3 * time.Second}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)

	if err != nil {
		p("error happened....", err)
	}

	defer res.Body.Close()
	if err != nil {
		fmt.Printf("[URL: %s] Faild\n", url)
		status = 0
		return
	}

	if res.StatusCode == 200 {
		fmt.Printf("GET[URL: %s] Success\n", url)
		body := res.Body
		bodyByte, _ := ioutil.ReadAll(body)
		resStr := string(bodyByte)
		content = resStr
		status = 200
		return
	} else {
		fmt.Printf("[URL: %s] Faild\n [STATUS: %d]\n", url, res.StatusCode)
	}
	return
}

func getUserAgent() string {
	var userAgent = [...]string{"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
		"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
		"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
		"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
		"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
		"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
		"Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
		"MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return userAgent[r.Intn(len(userAgent))]
}

func downWorker(uCh chan string) <-chan string {
	out := make(chan string, 10)
	go func() {
		for !done {
			<-time.After(time.Second)
			p("downloadworker wating for url")
			url := <-uCh
			p("download get url")
			if visted[url] == true {
				p("Found duplicated in: ", url)
				continue
			}
			visted[url] = true
			cont, status := httpGet(url)
			if status == 200 {
				// p(i1, l1)
				p("put the content ret")
				out <- cont
				p("put done")
				// fmt.Println("send content success ====>")
			} else {
				fmt.Printf("[URL: %s] Faild", url)
				fmt.Println("T_T Something Wrong...")
			}
		}
	}()
	return out
}

func parseWorker(uCh chan string, cCh <-chan string) <-chan string {
	out := make(chan string, 10)

	go func() {
		for !done {
			p("parserworker wating for content....")
			p("len(uChan): ", len(uCh))
			p("parserworker get content")
			cont := <-cCh
			magLis, urlLis, err := parse(cont)
			if err != nil {
				fmt.Println("errror on parse")
			}
			go func(urllis []string) {
				for _, url := range urllis {
					u := host + url[6:]
					if visted[u] == true {
						p("remove duplicated url")
						continue
					}
					if len(uCh) == cap(uCh) {
						return
					}
					p("put url in urlChan.....")
					uCh <- u
					p("put url in urlChan.....")
				}
			}(urlLis)
			if len(magLis) > 0 {
				go func(maglis []string) {
					for _, m := range maglis {
						if m != "" {
							p("put magnet in out chan==>.....")
							out <- m
							fmt.Println("magnet send...")
						}
					}
				}(magLis)
			}
		}
	}()
	return out
}

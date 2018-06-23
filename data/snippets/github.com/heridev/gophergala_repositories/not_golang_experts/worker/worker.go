package worker

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gophergala/not_golang_experts/model"
	"github.com/gophergala/not_golang_experts/notificator"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
	"log"
)

var stopchannel chan bool
var ticker *time.Ticker

func StartObserving(stopped chan bool) {
	stopchannel = stopped

	ticker = time.NewTicker(time.Millisecond * 120000) // 2 min

	go observe()
}

func StopObserving() {
	ticker.Stop()
	stopchannel <- true
}

func observe() {
	for t := range ticker.C {
		pagestocheck := model.PagesToCheck()
		for _, page := range pagestocheck {
			log.Printf("Checking page: %v - %v\n", page.Url, t)

			resultchan := make(chan string, 1)
			go requestHTML(page, resultchan)
			resultString := <-resultchan

			if page.HtmlString != resultString {
				page.HtmlString = resultString
				notificator.SendPageUpdatedNotificationToUsers(page.SubscribedUsersEmails(), page.Url)
				log.Println("UPDATED -> " + resultString)
			} else {
				page.LastCheckedAt = time.Now()
			}
			page.Save()
		}
	}
}

func requestHTML(p *model.Page, resultchan chan string) {
	res, err := http.Get(p.Url)

	if err != nil {
		panic(err)
	} else {
		defer res.Body.Close()
		html, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		re := regexp.MustCompile("<html(\\S|\\s)*\\/html>")
		matches := re.FindString(string(html))

		hasher := md5.New()
		hasher.Write([]byte(string(matches)))
		resultchan <- hex.EncodeToString(hasher.Sum(nil))
	}
}

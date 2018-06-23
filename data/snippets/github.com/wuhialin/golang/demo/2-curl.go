package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func main() {
	i := 0
	maxPage := 1 << 31
	reg := regexp.MustCompile(`<span\s+id="gamesPage".*?>\s*(\d+)\s*</span>`)
	for {
		i++
		if i > maxPage {
			break
		}
		url := "http://www.wowpower.com/showNewGame?page=" + strconv.Itoa(i)
		response, err := http.Get(url)
		if err != nil {
			continue
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			continue
		}
		q := reg.FindStringSubmatch(string(body))
		if q != nil {
			maxPage, _ = strconv.Atoi(q[1])
		}
		log.Println(i)
	}
}

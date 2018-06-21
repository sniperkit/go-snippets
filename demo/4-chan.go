/*
* @Author: wuhailin
* @Date:   2017-10-24 15:23:53
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-10-24 15:26:22
 */

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var baseUrl string = `http://www.wowpower.com/showNewGame?page=`

func main() {
	i := 1
	c := http.Client{}
	c.Timeout = 3 * time.Second
	req := regexp.MustCompile(`<span\s+id="gamesPage".*?>\s*(\d+)\s*</span>`)
	url := "http://www.wowpower.com/showNewGame?page=" + strconv.Itoa(i)
	response, err := c.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		q := req.FindStringSubmatch(string(body))
		formatter(body)
		log.Println(q)
	}
}

func formatter(body []byte) []map[string]string {
	games := []map[string]string{}
	body = regexp.MustCompile(`\n`).ReplaceAll(body, []byte(""))
	body = regexp.MustCompile(`\s+`).ReplaceAll(body, []byte(" "))

	nameReg := regexp.MustCompile(`<dd\s+class="game-name">(.+?)</dd>`)
	for _, q := range nameReg.FindAllSubmatch(body, -1) {
		tmpMap := map[string]string{}
		tmpMap["name"] = string(q[1])
		games = append(games, tmpMap)
	}

	priceReg := regexp.MustCompile(`<span>押金：(\d+)纳币</span>`)
	for k, q := range priceReg.FindAllSubmatch(body, -1) {
		games[k]["price"] = string(q[1])
	}

	dayReg := regexp.MustCompile(`<span>周期：(\d+)天</span>`)
	for k, q := range dayReg.FindAllSubmatch(body, -1) {
		games[k]["day"] = string(q[1])
	}

	gainProfitReg := regexp.MustCompile(`<span>返利：(\d+(\.\d+)?)纳币</span>`)
	for k, q := range gainProfitReg.FindAllSubmatch(body, -1) {
		games[k]["gain_profit"] = string(q[1])
	}

	return games
}

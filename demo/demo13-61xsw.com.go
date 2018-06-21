/*
* @Author: wuhailin
* @Date:   2017-12-01 14:11:21
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-12-01 16:40:07
 */
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var db *sql.DB

func main() {
	start := time.Now()
	mainUrl := "http://m.61xsw.com"
	response, err := http.Get(mainUrl)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	bodyString := string(body)
	bodyString = strings.Replace(bodyString, "\n", "", -1)
	exp := regexp.MustCompile(`<a\s+.*?href=['"](/?book-bid-\d+\.html?).*?['"].*?>(.*?)<\/a>`)
	if err != nil {
		log.Fatalln(err)
	}
	db, err = sql.Open("mysql", `root:123456@/yii?charset=utf8`)
	defer db.Close()
	if err != nil {
		log.Println(err)
	}
	urls := make(map[string]string)
	var pUrl *url.URL
	var insertSqlPlaceholder []string
	var insertSqlParams []interface{}
	for _, str := range exp.FindAllStringSubmatch(bodyString, -1) {
		if str[1] == "" {
			continue
		}
		pUrl, err = url.Parse(str[1])
		if err != nil {
			continue
		}
		if pUrl.String() == "" || pUrl.String() == "/" || urls[pUrl.String()] != "" {
			continue
		}
		urls[pUrl.String()] = pUrl.String()
		insertSqlPlaceholder = append(insertSqlPlaceholder, `(?, ?, ?)`)

		insertSqlParams = append(insertSqlParams, pUrl.String())
		insertSqlParams = append(insertSqlParams, mainUrl)
		insertSqlParams = append(insertSqlParams, fmt.Sprintf("%d", start.Unix()))
	}
	if len(insertSqlPlaceholder) > 0 {
		sqls := []string{"INSERT INTO crawl_txt (path, domain, created_at) VALUE"}
		sqls = append(sqls, strings.Join(insertSqlPlaceholder, ", "))
		sqls = append(sqls, `ON DUPLICATE KEY UPDATE path=VALUES(path)`)
		querySql := strings.Join(sqls, " ")
		_, err = db.Exec(querySql, insertSqlParams...)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println(time.Since(start))
}

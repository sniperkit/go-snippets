package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	//"regexp"
)

func main() {
	queryUrl := "http://wuhialin.free.ngrok.cc/vpn"
	data := []string{"http://www.170mv.com/"}
	response, err := http.PostForm(queryUrl, url.Values{"data": data})
	if err != nil {
		log.Fatalln(err)
	}
	if response.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		println(string(body))
	}
	//content, _ := ioutil.ReadFile("D:/index.html")
	//downloadUrlExp := regexp.MustCompile(`<a\s+id="video_down"\s+href="(.+)?"\s+title`)
	//matches := downloadUrlExp.FindSubmatch(content)
	//nameExp := regexp.MustCompile(`<title>(.*)?-高清MV下载地址-.*?</title>`)
	//nameMatches := nameExp.FindSubmatch(content)
	//log.Println(string(matches[1]), string(nameMatches[1]))
}

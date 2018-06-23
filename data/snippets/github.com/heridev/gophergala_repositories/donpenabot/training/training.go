package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

func main() {
	anaconda.SetConsumerKey(os.Getenv("APIKEY"))
	anaconda.SetConsumerSecret(os.Getenv("APISECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("ACCESSTOKEN"), os.Getenv("ACCESSTOKENSECRET"))
	v := url.Values{}
	v.Set("count", "1000")
	searchResult, err := api.GetSearch("epn", v)

	if err != nil {
		panic(err)
	}

	f, err := os.Create("seeds.txt")
	defer f.Close()
	buffer := bytes.Buffer{}
	lastTweet := ""
	for _, tweet := range searchResult.Statuses {
		if lastTweet != tweet.Text {
			hashtags := ""
			urls := ""
			if len(tweet.Entities.Hashtags) > 0 {
				for _, hash := range tweet.Entities.Hashtags {
					hashtags = (hashtags + hash.Text + " ")
				}
			} else {
				hashtags = "none"
			}
			if len(tweet.Entities.Urls) > 0 {
				for _, url := range tweet.Entities.Urls {
					fmt.Println(url.Display_url)
					urls = (urls + url.Display_url + " ")
				}
			} else {
				urls = "none"
			}
			text := strings.Replace(tweet.Text, "\"", "\\\"", -1)
			text = strings.Replace(text, "\n", "", -1)
			tweetbytes := []byte(fmt.Sprintf("\"%s\",\"%s\",\"%s\"\n", text, strings.TrimSpace(hashtags), urls))
			buffer.Write(tweetbytes)
			lastTweet = tweet.Text
		} else {
			lastTweet = tweet.Text
		}
	}
	f.Write(buffer.Bytes())
	f.Sync()
}

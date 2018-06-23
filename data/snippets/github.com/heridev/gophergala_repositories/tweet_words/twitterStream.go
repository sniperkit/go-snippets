package tweet_words

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/robfig/config"
	"gopkg.in/mgo.v2"
	"net/url"
	"strings"
	"time"
)

var TwitterApi *anaconda.TwitterApi
var Conf map[string]string

type TweetStore struct {
	TwitterURL     string
	Tweet          string
	Classification string
}

func init() {
	myconf, err := config.ReadDefault("app.properties")
	if err != nil {
		panic(err)
	}
	Conf = make(map[string]string)
	var keys = []string{"CONSUMER_KEY", "CONSUMER_SECRET", "ACCESS_TOKEN", "ACCESS_TOKEN_SECRET", "MONGO"}
	for key := range keys {
		Conf[keys[key]], err = myconf.String("", keys[key])
		if err != nil {
			panic(err)
		}
	}
	anaconda.SetConsumerKey(Conf["CONSUMER_KEY"])
	anaconda.SetConsumerSecret(Conf["CONSUMER_SECRET"])
	TwitterApi = anaconda.NewTwitterApi(Conf["ACCESS_TOKEN"], Conf["ACCESS_TOKEN_SECRET"])
}

func Tweets(query url.Values, timeout time.Duration, quit chan bool) <-chan anaconda.Tweet {
	stream, err := TwitterApi.UserStream(query) //PublicStreamFilter ?
	if err != nil {
		panic(err)
	}
	var tweet anaconda.Tweet
	var junk interface{}
	tweetChan := make(chan anaconda.Tweet, 1024) // as much as I like the consumer routine...
	go func() {
		quitter := time.After(timeout)
		for {
			select {
			case junk = <-stream.C:
				switch junk.(type) {
				case anaconda.Tweet:
					tweet = junk.(anaconda.Tweet)
					tweetChan <- tweet
				}
			case <-quitter:
				quit <- true
				return
			}
		}
	}()
	return tweetChan
}

func StoreTweets(query url.Values, timeout time.Duration, collectionName string) (retChan chan bool) {
	retChan = make(chan bool)
	quit := make(chan bool)
	tweetsChan := Tweets(query, timeout, quit)
	go func() {
		mgoSession, err := mgo.Dial(Conf["MONGO"])
		if err != nil {
			panic(err)
		}
		mgoSession.SetMode(mgo.Monotonic, true)
		defer mgoSession.Close()
		var twitterUrl, classification string
		for {
			select {
			case tweet := <-tweetsChan:
				newSes := mgoSession.Copy()
				defer newSes.Close()
				col := newSes.DB("test").C(collectionName)
				if col == nil {
					panic("unable to get collection")
				}
				twitterUrl = fmt.Sprintf("https://www.twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IdStr)
				tweet.Text = strings.Replace(tweet.Text, "\n", "", -1)
				classification = ClassifyTweet(tweet.Text)
				err = col.Insert(&TweetStore{twitterUrl, tweet.Text, classification})
				if err != nil {
					panic(err)
				}
			case <-quit:
				retChan <- true
				return
			}
		}
	}()
	return
}

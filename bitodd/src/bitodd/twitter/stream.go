package twitter

import (
	"bitodd/config"
	"bitodd/model"
	"errors"
	"github.com/darkhelmet/twitterstream"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func loadCredentials() (client *twitterstream.Client, err error) {
	credentials, err := ioutil.ReadFile("CREDENTIALS")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(credentials), "\n")
	if len(lines) < 4 {
		return nil, errors.New("Credentials file is too short (lines < 4)")
	}

	return twitterstream.NewClient(lines[0], lines[1], lines[2], lines[3]), nil
}

var maxWait = 600
var wait = 15

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func OpenStream() {
	client, err := loadCredentials()
	if err != nil {
		log.Fatalln("Error while loading credentials", err)
	}

	for {
		keywords := config.GetConfig().Keywords
		log.Printf("Started streaming keywords: %s", keywords)
		conn, err := client.Track(keywords...)
		if err != nil {
			log.Printf("Tracking failed: %s", err)
			wait = wait << 1
			log.Printf("Waiting for %d seconds before reconnect", min(wait, maxWait))
			time.Sleep(time.Duration(min(wait, maxWait)) * time.Second)
			continue
		} else {
			wait = 15
		}
		decodeTweets(conn)
	}
}

func decodeTweets(conn *twitterstream.Connection) {
	for {
		if tweet, err := conn.Next(); err == nil {

			// Use only tweets with coordinates
			if tweet.Coordinates != nil {
				model.PostLocationTweet(tweet)
			} else {
			}

		} else {
			log.Printf("Decoding tweet failed: %s", err)
			conn.Close()
			return
		}
	}
}

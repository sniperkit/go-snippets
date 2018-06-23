package reddit

import (
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"strconv"
	"time"
)

type QueryResponse struct {
	Data struct {
		Children []struct {
			Data Post `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type Post struct {
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	RedditLink string    `json:"permalink"`
	Retrieved  time.Time `json:"-"`
}

func QueryRedditTop(c context.Context, sub string, limit int, t string) ([]Post, error) {
	client := urlfetch.Client(c)
	url := "https://www.reddit.com/r/" + sub + "/top/.json?t=" + t + "&limit=" + strconv.Itoa(limit)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var parsed QueryResponse
	if err = decoder.Decode(&parsed); err != nil {
		return nil, err
	}
	var posts []Post
	for _, element := range parsed.Data.Children {
		post := element.Data
		post.RedditLink = "https://reddit.com" + post.RedditLink
		post.Retrieved = time.Now()
		posts = append(posts, post)
	}
	return posts, nil
}

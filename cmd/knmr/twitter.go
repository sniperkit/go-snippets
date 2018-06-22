package main

import (
	"fmt"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Ken2mer/knmr/logger"
	"github.com/urfave/cli"
)

var (
	placeID         int64  = 1118550    // Yokohama
	slug            = "engineer" // for GetListTweetsBySlug()
	ownerScreenName = "Ken2mer"  // for GetListTweetsBySlug()
	queryString     = "lang:ja"  // for GetSearch()
)

type twClient struct {
	api *anaconda.TwitterApi
}

var twitterCommand = cli.Command{
	Name:   "twitter",
	Usage:  "twitter",
	Action: twitterCmd,
}

func twitterCmd(ctx *cli.Context) error {
	c := newTwitterClient()
	c.dumpTrendResp()
	return nil
}

func newTwitterClient() twClient {
	// consumer_key        string = "XXXXXX"
	// consumer_secret     string = "XXXXXX"
	// access_token        string = "XXXXXX"
	// access_token_secret string = "XXXXXX"
	anaconda.SetConsumerKey(consumer_key)
	anaconda.SetConsumerSecret(consumer_secret)
	api := anaconda.NewTwitterApi(access_token, access_token_secret)
	return twClient{api: api}
}

func (c *twClient) dumpActivity() {
	url, err := c.api.GetActivityWebhooks(nil)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	logger.DumpJSON(url)
}

func (c *twClient) dumpTrendResp() {
	trendResp, err := c.api.GetTrendsByPlace(placeID, nil)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	logger.DumpJSON(trendResp)
}

func (c *twClient) dumpListTweets() {
	tweets, err := c.api.GetListTweetsBySlug(slug, ownerScreenName, true, nil)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	logger.DumpJSON(tweets)
}

func (c *twClient) dmupTimeline() {
	timeline, err := c.api.GetHomeTimeline(nil)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	logger.DumpJSON(timeline)
}

func (c *twClient) dumpFriends() {
	cursor, err := c.api.GetFriendsList(nil)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	for i, user := range cursor.Users {
		fmt.Printf("%d: \n", i)
		fmt.Printf("user: %s\n", user.Name)
	}
}

func (c *twClient) dumpFollowers() {
	cursor, err := c.api.GetFollowersList(nil)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	for i, user := range cursor.Users {
		fmt.Printf("%d: \n", i)
		fmt.Printf("user: %s\n", user.Name)
	}
}

func (c *twClient) dumpFavorites() {
	favorites, err := c.api.GetFavorites(nil)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	for i, tweet := range favorites {
		fmt.Printf("%d: \n", i)
		fmt.Printf("%s wrote: \n", tweet.User.Name)
		fmt.Printf("%s\n\n", tweet.FullText)
	}
}

func (c *twClient) dumpSearchResult() {
	api := c.api
	searchResponse, err := c.api.GetSearch(queryString, nil)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	for searchResponse.Statuses != nil {
		logger.DumpJSON(searchResponse)
		searchResponse, err = searchResponse.GetNext(api)
	}
}

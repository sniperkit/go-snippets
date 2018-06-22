package main

/*
*	The purpose of Collector.go is to abstract the collection
* process. Collector.go interfaces with Controller.go that
* send the commands to Collector.go to initate raw Twitter
* data collection. This data is then sent back to Collector.go
* and sent to Formatter.go which correctly formats the data for
* consistency and clairty. The formatted data comes back to
* Collector.go and handed off to Controller.go.
* Be careful about routing the Twitter data as the raw data will
* cause disruptions with the rest of Clay unless properly formatted.
 */

type Collector struct { // needs to be struct so I can call funcs from it
	twComm *twitterComm
	f      formatter
}

func newUserClientCollector() Collector {
	return Collector{NewTwitterUserClient(), newFormatter()}
}

func newAppClientCollector() Collector {
	return Collector{NewTwitterAppClient(), newFormatter()}
}

/*
func (c Collector) collectFullInfo(userID int64) (map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {
	u := c.CollectUser(userID)
	tweets := c.collectTweets(u["id"].(int64), 200)
	favs := c.collectFavourites(u["id"].(int64), 200)
	for _, fav := range favs {
		u["favourite_ids"] = append(u["favourite_ids"].([]int64), c.f.FormatToInt64(fav["id"]))
	}
	u = c.f.FormatUser(u) //reformat user after collecting favourites
	return u, tweets, favs
}*/

func (c Collector) CollectFavourites(userID int64) []map[string]interface{} {
	favs := c.twComm.Favourites(userID)
	for _, fav := range favs {
		c.f.FormatTweet(fav)
	}
	return favs
}

func (c Collector) CollectFollowing(userID int64) []int64 {
	unformattedIDS := c.twComm.FollowingIDS(userID)
	return c.f.FormatArray(unformattedIDS)
}

func (c Collector) CollectFollowers(userID int64) []int64 {
	unformattedIDS := c.twComm.FollowerIDS(userID)
	return c.f.FormatArray(unformattedIDS)
}

func (c Collector) CollectUsers(userIDS []int64) []map[string]interface{} {
	users := c.twComm.Users(userIDS)
	for _, user := range users {
		c.f.FormatUser(user)
	}
	return users //have been formatted
}

func (c Collector) CollectUser(userID int64) map[string]interface{} {
	user := map[string]interface{}{}
	for len(user) == 0 {
		user = c.twComm.User(userID)
	}
	return c.f.FormatUser(user)
}

func (c Collector) CollectTweets(userID int64) []map[string]interface{} {
	tweets := c.twComm.Tweets(userID)
	for _, tw := range tweets {
		c.f.FormatTweet(tw)
	}
	return tweets
}

/*
func (c Collector) collectFollowers(user *map[string]interface{}, count int) []map[string]interface{} {
	rawFollowers := c.twComm.followers(user, count)
	formattedFollowers := []map[string]interface{}{}
	for _, fol := range *rawFollowers {
		formattedFollowers = append(formattedFollowers, c.f.FormatUser(fol))
	}
	return formattedFollowers
}*/
/*
func (c Collector) collectFollowing(user *map[string]interface{}, count int) []map[string]interface{} {
	rawFollowing := c.twComm.following(user, count)
	formatteFollowing := []map[string]interface{}{}
	for _, fol := range *rawFollowing {
		formatteFollowing = append(formatteFollowing, c.f.FormatUser(fol))
	}
	return formatteFollowing
}*/

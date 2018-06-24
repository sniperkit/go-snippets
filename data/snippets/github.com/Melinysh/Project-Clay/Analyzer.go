package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/deckarep/golang-set"
)

/*
*  This is to be used to analyze Twitter data
* that has already been collected and formatted.
* It is responsible for detecting changes, for example,
* new followers/following, and delegate the event creation
* to EventCreator.go. The event objects returned by
* EventCreator.go are then passed onto MongoComm.go
* for saving along with all other relevant data. The
* analyzer also requests further collection relating to
* the changes detected in Twitter data, for example,
* if a new follower is detected, fetch that follower and
* all relevant data.
 */

type Analyzer struct {
	m mongoComm
	e EventCreator
}

func newAnalyzer() Analyzer {
	return Analyzer{NewMongoComm() /* NewRethinkComm()*/, newEventCreator()}
}

func (a Analyzer) AnalyzeTweets(tweets []map[string]interface{}) {
	for _, tw := range tweets {
		a.analyzeTweet(tw)
	}
}

func (a Analyzer) AnalyzeFavouritesForUser(favs []map[string]interface{}, userID int64) {
	for _, fav := range favs {
		a.analyzeTweet(fav)
	}

	if user, err := a.m.User(userID); err != nil {
		fmt.Println("User not found in DB. Just saving tweets for now.")
		return
	} else {
		favIDS := []int64{} //Set new ids into array and get diffs
		for _, fav := range favs {
			favIDS = append(favIDS, fav["id"].(int64))
		}
		newFavs := diff(favIDS, convertInterfaces(user["favourite_ids"].([]interface{})))
		//		fmt.Println("ANALYZER: New favs for user", user["screen_name"], newFavs)
		a.newEvents(user["id"].(int64), newFavs, a.e.newFavouriteEvent)
		user["favourite_ids"] = append(convertInterfaces(user["favourite_ids"].([]interface{})), newFavs...) //existing plus new
		a.m.SaveUser(user)                                                                                   //makes sure those new fav ids are saved
	}
}

func (a Analyzer) analyzeTweet(tweet map[string]interface{}) {

	savedTweet, err := a.m.Tweet(tweet["id"].(int64))

	if err != nil {
		a.m.SaveTweet(tweet)
		return
	}

	if changes, hasChanged := a.detectTweetdifferences(tweet, savedTweet); hasChanged {
		a.m.SaveTweet(changes)
		//		changedEvent := a.e.newChangeEvent(tweet["id"].(int64), changes)
		//		a.m.SaveEvent(changedEvent)
	}
}

func (a Analyzer) AnalyzeFollowing(ids []int64, userID int64) {
	savedUser, err := a.m.User(userID)

	if err != nil {
		fmt.Println("Error:", err, "Fetched following ids for a user who threw error so Analyzer can't proceed with them. Discarding...")
		return
	}

	newFollowing := diff(ids, convertInterfaces(savedUser["following_ids"].([]interface{})))
	a.newEvents(userID, newFollowing, a.e.newFollowingEvent)

	newUnfollowing := diff(convertInterfaces(savedUser["following_ids"].([]interface{})), ids)
	a.newEvents(userID, newUnfollowing, a.e.newUnfollowingEvent)
	a.m.SaveUser(map[string]interface{}{"id": userID, "following_ids": ids})
	fmt.Println("Completed analysis of following for", userID, newFollowing, newUnfollowing)
}

func (a Analyzer) AnalyzeFollowers(ids []int64, userID int64) {
	savedUser, err := a.m.User(userID)

	if err != nil {
		fmt.Println("Error:", err, "Fetched follower ids for a user who threw error so Analyzer can't proceed with them. Discarding...")
		return
	}
	newFollowers := diff(ids, convertInterfaces(savedUser["follower_ids"].([]interface{})))
	a.newEvents(userID, newFollowers, a.e.newFollowerEvent)

	newUnfollowers := diff(convertInterfaces(savedUser["follower_ids"].([]interface{})), ids)
	a.newEvents(userID, newUnfollowers, a.e.newUnfollowerEvent)
	a.m.SaveUser(map[string]interface{}{"id": userID, "follower_ids": ids})
	fmt.Println("Completed analysis of followers for", userID, newFollowers, newUnfollowers)
}

func (a Analyzer) AnalyzeUser(user map[string]interface{}) {
	savedUser, err := a.m.User(user["id"].(int64))

	if err != nil {
		a.m.SaveUser(user)
		return
	}

	if status, ok := user["status"]; ok {
		a.analyzeTweet(status.(map[string]interface{}))
	}

	if changes, hasChanged := a.detectUserdifferences(user, savedUser); hasChanged {
		//		fmt.Println("Here's whats changed in the user", (*user)["id"], *changes)
		defer a.m.SaveUser(changes) // so it gets called by both return and the function finishing if and only if there are changes.

		if len(user["following_ids"].([]int64)) == 0 && len(user["follower_ids"].([]int64)) == 0 {
			//fmt.Println("Basic user", user["screen_name"], "completed.")
			delete(changes, "following_ids")
			delete(changes, "follower_ids")
			return
		}

		newFollowers := diff((user["follower_ids"]).([]int64), convertInterfaces(savedUser["follower_ids"].([]interface{})))
		a.newEvents(user["id"].(int64), newFollowers, a.e.newFollowerEvent)

		newUnfollowers := diff(convertInterfaces(savedUser["follower_ids"].([]interface{})), (user["follower_ids"]).([]int64))
		a.newEvents(user["id"].(int64), newUnfollowers, a.e.newUnfollowerEvent)

		newFollowing := diff((user["following_ids"]).([]int64), convertInterfaces(savedUser["following_ids"].([]interface{})))
		a.newEvents(user["id"].(int64), newFollowing, a.e.newFollowingEvent)

		newUnfollowing := diff(convertInterfaces(savedUser["following_ids"].([]interface{})), (user["following_ids"]).([]int64))
		a.newEvents(user["id"].(int64), newUnfollowing, a.e.newUnfollowingEvent)
		//	fmt.Println("@@ Test dump of events", *newFollowers, *newUnfollowers, *newFollowing, *newUnfollowing)
		//save what has changed
		/*		changedEvent := a.e.newChangeEvent(user["id"].(int64), changes)
				fmt.Println("DEBUG: Changes for user", user["id"].(int64), changedEvent)
				a.m.SaveEvent(changedEvent)*/
		fmt.Println("Updated user", user["screen_name"], ":", user["id"])
	}
}

//Checks for user diffs and guarantees that the user will be in the db upon return
func (a Analyzer) detectUserdifferences(newUser map[string]interface{}, savedUser map[string]interface{}) (map[string]interface{}, bool) {

	changesMap := map[string]interface{}{"id": savedUser["id"], "lastModified": time.Now().Unix()}
	for key, newValue := range newUser {
		if key == "favourite_ids" && len(newValue.([]int64)) == 0 {
			continue //The user will usually have 0 fav ids until it goes through the favs queue, so ignore this key.
		}
		if savedValue, ok := savedUser[key]; ok {
			if reflect.DeepEqual(newValue, savedValue) == false {
				changesMap[key] = newValue
			}
		} else { //no key in saved map
			changesMap[key] = newValue //key not currently saved so save it
		}
	}

	if len(changesMap) > 2 { //The id and lastMod date is always include in the map by default
		return changesMap, true
	}
	fmt.Println("No changes detected for", newUser["id"], newUser["screen_name"])
	return changesMap, false //they are the exact same
}

func (a Analyzer) detectTweetdifferences(newTweet map[string]interface{}, savedTweet map[string]interface{}) (map[string]interface{}, bool) {

	changesMap := map[string]interface{}{"id": savedTweet["id"], "lastModified": time.Now().Unix()}
	for key, newValue := range newTweet {
		if savedValue, ok := savedTweet[key]; ok {
			if reflect.DeepEqual(newValue, savedValue) == false {
				changesMap[key] = newValue
			}
		} else {
			changesMap[key] = newValue //key not currently saved so save it
		}
	}

	if len(changesMap) > 2 { //The id and lastMod date is always include in the map by default
		return changesMap, true
	}
	return changesMap, false //they are the exact same
}

func (a Analyzer) newEvents(i int64, ids []int64, eventFn func(int64, int64) map[string]interface{}) {
	for _, id := range ids {
		e := eventFn(i, id)
		a.m.SaveEvent(e)
	}
}

// Pass in new set then exsisting set for latest additions.
func diff(a1 []int64, a2 []int64) []int64 {
	s1 := mapset.NewSetFromSlice(convertInts(a1))
	s2 := mapset.NewSetFromSlice(convertInts(a2))
	diff := s1.Difference(s2)
	if len(diff.ToSlice()) == 0 {
		return []int64{}
	}
	return convertInterfaces(diff.ToSlice())
}

func convertInts(a []int64) []interface{} {
	s := make([]interface{}, len(a))
	for i, v := range a {
		s[i] = v
	}
	return s
}

func convertInterfaces(a []interface{}) []int64 {
	s := make([]int64, len(a))
	for i, v := range a {
		s[i] = v.(int64)
	}
	return s
}

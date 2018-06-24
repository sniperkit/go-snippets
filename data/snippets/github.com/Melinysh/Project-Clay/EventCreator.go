package main

import "time" //for timestamps

/*
*	EventCreator.go is a simple module that recieves information
* about an event, such as type and the IDS of items involved, and
* creates an event map[string]interface{} with it. This map is subsequently
* returned to Analyzer.go.
 */

type EventCreator struct {
	blah int // needed for I can call funcs from it
}

func newEventCreator() EventCreator {
	return EventCreator{}
}

func (e EventCreator) newFavouriteEvent(userID int64, tweetID int64) map[string]interface{} {
	event := map[string]interface{}{}
	event["id"] = userID
	event["type"] = "favourited"
	event["timestamp"] = time.Now().Unix()
	event["tweetID"] = tweetID
	return event
}

func (e EventCreator) newUnfavouriteEvent(userID int64, tweetID int64) map[string]interface{} {
	event := map[string]interface{}{}
	event["id"] = userID
	event["type"] = "unfavourited"
	event["timestamp"] = time.Now().Unix()
	event["tweetID"] = tweetID
	return event
}

func (e EventCreator) newFollowerEvent(userID int64, folID int64) map[string]interface{} {
	event := map[string]interface{}{}
	event["id"] = userID
	event["type"] = "followedBy"
	event["timestamp"] = time.Now().Unix()
	event["userID"] = folID
	return event
}

func (e EventCreator) newUnfollowerEvent(userID int64, unfolID int64) map[string]interface{} {
	event := map[string]interface{}{}
	event["id"] = userID
	event["type"] = "unfollowedBy"
	event["timestamp"] = time.Now().Unix()
	event["userID"] = unfolID
	return event
}

func (e EventCreator) newFollowingEvent(userID int64, folID int64) map[string]interface{} {
	event := map[string]interface{}{}
	event["id"] = userID
	event["type"] = "followed"
	event["timestamp"] = time.Now().Unix()
	event["userID"] = folID
	return event
}

func (e EventCreator) newUnfollowingEvent(userID int64, unfolID int64) map[string]interface{} {
	event := map[string]interface{}{}
	event["id"] = userID
	event["type"] = "unfollowed"
	event["timestamp"] = time.Now().Unix()
	event["userID"] = unfolID
	return event
}

func (e EventCreator) newChangeEvent(id int64, changes map[string]interface{}) map[string]interface{} {
	event := map[string]interface{}{}
	event["id"] = id
	event["type"] = "changed"
	event["timestamp"] = time.Now().Unix()
	//remove 'lastModified in changesMap as it is rendundant with timestamp
	delete(changes, "lastModified")
	event["changes"] = changes
	return event
}

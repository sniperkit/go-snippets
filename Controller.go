package main

import (
	"fmt"
	"time"

	"github.com/oleiade/lane"
)

type Controller struct {
	a Analyzer
	c Collector
}

var myAnalyzer = newAnalyzer()

var usersInPrimaryQueue = &map[int64]bool{}
var usersInSecondaryQueue = &map[int64]bool{}
var usersInTertiaryQueue = &map[int64]bool{}
var usersInRelationsQueue = &map[int64]bool{}

var globalPrimaryUserFetchQueue = lane.NewQueue()
var globalSecondaryUserFetchQueue = lane.NewQueue()
var globalTertiaryUserFetchQueue = lane.NewQueue()

var globalAppFavouritesFetchQueue = lane.NewQueue()
var globalAppTweetFetchQueue = lane.NewQueue()

var globalUserTweetFetchQueue = lane.NewQueue()
var globalUserFavouritesFetchQueue = lane.NewQueue()

var globalRelationsFetchQueue = lane.NewQueue()

func main() {
	fmt.Println("\n\tWelcome.")
	fmt.Println("-=-=-=- CLAY v1.2.1 =-=-=-=")
	fmt.Println("Let's begin.")

	go processPrimaryUserQueue()
	go processSecondaryUserQueue()
	go processTertiaryUserQueue()
	go processAppTweetQueue()
	go processUserTweetQueue()
	go processAppFavQueue()
	go processUserFavQueue()
	go processRelationsQueue()

	start() //should never return.
}

func start() {
	id := int64(0) // My Twitter ID
	myController := Controller{myAnalyzer, newUserClientCollector()}
	for {
		user := myController.c.CollectUser(id)
		globalUserFavouritesFetchQueue.Enqueue(id)
		globalUserTweetFetchQueue.Enqueue(id)
		go myController.a.AnalyzeUser(user)
		following := myController.c.CollectFollowing(id)
		followers := myController.c.CollectFollowers(id)
		go myController.a.AnalyzeFollowing(following, id)
		go myController.a.AnalyzeFollowers(followers, id)
		for _, userID := range followers {
			if _, ok := (*usersInPrimaryQueue)[userID]; ok == false { //they aren't in the Queue or have been processed, so add them.
				globalPrimaryUserFetchQueue.Enqueue(userID)
			}
		}

		for _, userID := range following {
			if _, ok := (*usersInPrimaryQueue)[userID]; ok == false { //they aren't in the Queue or have been processed, so add them.
				globalPrimaryUserFetchQueue.Enqueue(userID)
			}
		}

		time.Sleep(time.Duration(60 * time.Minute))
	}
}

/*
func firstDegree(followers []int64, following []int64) {
	fmt.Println("Entering first degree.")
	for _, userID := range followers {
		if _, ok := (*usersInPrimaryQueue)[userID]; ok == false { //they aren't in the Queue or have been processed, so add them.
			globalPrimaryUserFetchQueue.Enqueue(userID)
		}
	}

	for _, userID := range following {
		if _, ok := (*usersInPrimaryQueue)[userID]; ok == false { //they aren't in the Queue or have been processed, so add them.
			globalPrimaryUserFetchQueue.Enqueue(userID)
		}
	}
	fmt.Println("Enqueued first degree.")
	processPrimaryUserQueue() // <- This will never return as dequeued objects are requeued
}*/

func secondDegree(followers []int64, following []int64) {
	fmt.Println("Entering a second degree.", followers, following)
	for _, userID := range followers {
		if _, ok := (*usersInSecondaryQueue)[userID]; ok == false { //they aren't in the Queue or have been processed, so add them.
			globalSecondaryUserFetchQueue.Enqueue(userID)
		}

	}

	for _, userID := range following {
		if _, ok := (*usersInSecondaryQueue)[userID]; ok == false { //they aren't in the Queue or have been processed, so add them.
			globalSecondaryUserFetchQueue.Enqueue(userID)
		}
	}
	fmt.Println("Enqueued more into second degree queue. Current capacity:", globalSecondaryUserFetchQueue.Deque.Size())
}

func thirdDegree(followers []int64, following []int64) {
	fmt.Println("Entering a third degree.")
	for _, userID := range followers {
		if _, ok := (*usersInTertiaryQueue)[userID]; ok == false { //they aren't in the Queue or have been processed, so add them.
			if _, secOk := (*usersInSecondaryQueue)[userID]; secOk == false { //not in 2nd or 3rd or 1st degree
				if _, primOk := (*usersInPrimaryQueue)[userID]; primOk == false {
					globalTertiaryUserFetchQueue.Enqueue(userID)
				}
			}
		}
	}

	for _, userID := range following {
		if _, ok := (*usersInTertiaryQueue)[userID]; ok == false { //they aren't in the Queue or have been processed, so add them.
			if _, secOk := (*usersInSecondaryQueue)[userID]; secOk == false { //not in 2nd or 3rd or 1st degree
				if _, primOk := (*usersInPrimaryQueue)[userID]; primOk == false {
					globalTertiaryUserFetchQueue.Enqueue(userID)
				}
			}
		}
	}
	fmt.Println("Enqueued more into the third degree queue. Current capacity:", globalTertiaryUserFetchQueue.Deque.Size())
}

func processAppTweetQueue() {
	myController := Controller{myAnalyzer, newAppClientCollector()}
	for {
		userID := globalAppTweetFetchQueue.Dequeue()
		if userID == nil {
			continue
		}
		tweets := myController.c.CollectTweets(userID.(int64))
		go myController.a.AnalyzeTweets(tweets)
	}
}

func processUserTweetQueue() {
	myController := Controller{myAnalyzer, newUserClientCollector()}
	for {
		userID := globalUserTweetFetchQueue.Dequeue()
		if userID == nil {
			continue
		}
		tweets := myController.c.CollectTweets(userID.(int64))
		myController.a.AnalyzeTweets(tweets)
	}
}

func processAppFavQueue() {
	myController := Controller{myAnalyzer, newAppClientCollector()}
	for {
		userID := globalAppFavouritesFetchQueue.Dequeue()
		if userID == nil {
			continue
		}
		favs := myController.c.CollectFavourites(userID.(int64))
		go myController.a.AnalyzeFavouritesForUser(favs, userID.(int64))
	}
}

func processUserFavQueue() {
	myController := Controller{myAnalyzer, newUserClientCollector()}
	for {
		userID := globalUserFavouritesFetchQueue.Dequeue()
		if userID == nil {
			continue
		}
		favs := myController.c.CollectFavourites(userID.(int64))
		myController.a.AnalyzeFavouritesForUser(favs, userID.(int64))
	}
}

func processPrimaryUserQueue() {
	myController := Controller{myAnalyzer, newUserClientCollector()}
	for {
		userID := globalPrimaryUserFetchQueue.Dequeue()
		if userID == nil {
			time.Sleep(time.Duration(5 * time.Second))
			continue
		}
		fmt.Println("Processing 1st degree user", userID)
		user := myController.c.CollectUser(userID.(int64))
		if user["shouldFetchFavourites"].(bool) == true {
			globalUserFavouritesFetchQueue.Enqueue(user["id"].(int64))
		}
		if user["shouldFetchTweets"].(bool) == true {
			globalUserTweetFetchQueue.Enqueue(user["id"].(int64))
		}
		user["degree"] = 1
		myController.a.AnalyzeUser(user)
		if user["shouldFetchFollowers"].(bool) == true && user["shouldFetchFollowing"].(bool) == true {
			globalRelationsFetchQueue.Enqueue([]interface{}{userID, secondDegree})
		}
		fmt.Println("Completed primary user", user["screen_name"], ":", user["id_str"])
	}
}

func processSecondaryUserQueue() {
	myController := Controller{myAnalyzer, newAppClientCollector()}
	for {
		userIDS := []int64{}
		for i := 0; i < 99; i++ {
			id := globalSecondaryUserFetchQueue.Dequeue()
			if id == nil {
				time.Sleep(time.Duration(5 * time.Second))
				continue
			}

			userIDS = append(userIDS, id.(int64))
		}
		users := myController.c.CollectUsers(userIDS)
		fmt.Println("Collected ", len(users), " for the secondary user queue.")
		for _, user := range users {
			//			fmt.Println("On", i, "user of secondary queue batch response.", user["screen_name"])
			if user["shouldFetchFavourites"].(bool) == true {
				globalAppFavouritesFetchQueue.Enqueue(user["id"].(int64))
			}
			if user["shouldFetchTweets"].(bool) == true {
				globalAppTweetFetchQueue.Enqueue(user["id"].(int64))
			}
			user["degree"] = 2
			myController.a.AnalyzeUser(user)
			if _, ok := (*usersInRelationsQueue)[(user["id"]).(int64)]; ok == false && user["shouldFetchFollowers"].(bool) == true && user["shouldFetchFollowing"].(bool) == true {
				globalRelationsFetchQueue.Enqueue([]interface{}{user["id"].(int64), thirdDegree})
			}
			delete((*usersInSecondaryQueue), (user["id"]).(int64))
			fmt.Println("Completed secondary user", user["screen_name"], ":", user["id_str"])
		}
	}
}

func processTertiaryUserQueue() {
	myController := Controller{myAnalyzer, newAppClientCollector()}
	for {
		userIDS := []int64{}
		for i := 0; i < 99; i++ {
			id := globalTertiaryUserFetchQueue.Dequeue()
			if id == nil {
				time.Sleep(time.Duration(5 * time.Second))
				continue
			}
			userIDS = append(userIDS, id.(int64))
		}

		users := myController.c.CollectUsers(userIDS)
		for _, user := range users {
			user["degree"] = 3
			myController.a.AnalyzeUser(user)
			delete((*usersInTertiaryQueue), (user["id"]).(int64))
			if _, ok := (*usersInRelationsQueue)[(user["id"]).(int64)]; ok == false && user["shouldFetchFollowers"].(bool) == true && user["shouldFetchFollowing"].(bool) == true {

				globalRelationsFetchQueue.Enqueue([]interface{}{user["id"].(int64)})
			}

			fmt.Println("Completed tertiary user", user["screen_name"], ":", user["id_str"])
		}
	}
}

func processRelationsQueue() {
	controller := Controller{myAnalyzer, newAppClientCollector()}
	for {
		out := globalRelationsFetchQueue.Dequeue() // index 0: id of user, index 1: func (second/third degree enqueue)
		if out == nil {
			time.Sleep(time.Duration(5 * time.Second))
			continue
		}
		relationRequestArray := out.([]interface{})
		id := relationRequestArray[0].(int64)
		followers := controller.c.CollectFollowers(id)
		controller.a.AnalyzeFollowers(followers, id)
		following := controller.c.CollectFollowing(id)
		controller.a.AnalyzeFollowing(following, id)
		if len(relationRequestArray) == 2 {
			function := relationRequestArray[1].(func([]int64, []int64))
			function(followers, following)
		}
		delete((*usersInRelationsQueue), id)
	}
}

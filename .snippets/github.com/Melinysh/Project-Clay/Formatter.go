package main

import (
	"fmt"
	"reflect"
	"strconv"
)

type formatter struct {
	blah int //need so I can call func (f formatter) from it
}

func newFormatter() formatter {
	return formatter{}
}

func (f formatter) FormatTweet(tw map[string]interface{}) map[string]interface{} {
	tw["id"] = f.FormatToInt64(tw["id"])
	tw["retweet_count"] = f.FormatToInt64(tw["retweet_count"])
	tw["favorite_count"] = f.FormatToInt64(tw["favorite_count"])
	if user, ok := tw["user"]; ok {
		u := f.FormatUser(user.(map[string]interface{}))
		tw["user"] = u
	}
	return tw
}

func (f formatter) FormatUser(user map[string]interface{}) map[string]interface{} {
	user["id"] = f.FormatToInt64(user["id"])
	user["followers_count"] = f.FormatToInt64(user["followers_count"])
	user["friends_count"] = f.FormatToInt64(user["friends_count"])
	user["favourites_count"] = f.FormatToInt64(user["favourites_count"])
	user["statuses_count"] = f.FormatToInt64(user["statuses_count"])

	if status, ok := user["status"]; ok {
		tw := f.FormatTweet(status.(map[string]interface{})) //doesn't create infinite loop b/c these status don't contain user object
		user["status"] = tw
	}

	if followingIDS, ok := user["following_ids"]; ok && reflect.TypeOf(user["following_ids"]).String() != "[]int64" {
		user["following_ids"] = f.FormatArray(followingIDS.([]interface{}))
	} else {
		user["following_ids"] = []int64{}
	}

	if followerIDS, ok := user["follower_ids"]; ok && reflect.TypeOf(user["follower_ids"]).String() != "[]int64" {
		user["follower_ids"] = f.FormatArray(followerIDS.([]interface{}))
	} else {
		user["follower_ids"] = []int64{}
	}

	if favouriteIDS, ok := user["favourite_ids"]; ok && reflect.TypeOf(user["favourite_ids"]).String() != "[]int64" {
		user["favourite_ids"] = f.FormatArray(favouriteIDS.([]interface{}))

	} else {
		user["favourite_ids"] = []int64{}
	}

	if user["followers_count"].(int64) > 1000 || user["friends_count"].(int64) > 1000 {
		user["shouldFetchFollowers"] = false
		user["shouldFetchFollowing"] = false
		user["shouldFetchTweets"] = false
		user["shouldFetchFavourites"] = false
	} else {
		user["shouldFetchFollowers"] = true
		user["shouldFetchFollowing"] = true
		user["shouldFetchTweets"] = true
		user["shouldFetchFavourites"] = true
	}
	return user
}

// FORMATING

func (f formatter) FormatArray(a []interface{}) []int64 {
	if len(a) == 0 {
		return []int64{}
	}
	int64Array := []int64{}
	for _, val := range a {
		int64Array = append(int64Array, f.FormatToInt64(val))
	}
	return int64Array
}

func (f formatter) FormatToInt64(val interface{}) int64 {
	var id int64
	var err error
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("RECOVERY in FormatToInt64 ", r, val)
			panic(r)
		}
	}()
	if reflect.TypeOf(val).Name() == "float64" {
		id = int64(val.(float64))
	} else if reflect.TypeOf(val).Name() == "int64" {
		id = val.(int64)
	} else if reflect.TypeOf(val).Name() == "string" {
		if id, err = strconv.ParseInt(val.(string), 10, 64); err != nil {
			fmt.Println("Error formatting string to int64", err)
		}
	} else if reflect.TypeOf(val).Name() == "int" {
		id = int64(val.(int))
	} else {
		fmt.Println("You passed", val, "which is of type", reflect.TypeOf(val).Name())
		panic("Non number value.")
	}
	return id
}

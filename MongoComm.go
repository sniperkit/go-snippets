package main

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoComm struct {
	session *mgo.Session
}

func NewMongoComm() mongoComm {
	var session *mgo.Session
	var err error
	if session, err = mgo.Dial("localhost:27017"); err != nil {
		fmt.Println("There was an error Dialing", err)
		panic("ERROR CONNECTING TO MONGO DB.")
	}
	return mongoComm{session}
}

func (mgoC mongoComm) SaveTweet(tw map[string]interface{}) {
	tweetColl := mgoC.session.DB("clayDB").C("Tweet")
	if _, err := tweetColl.Upsert(bson.M{"id": tw["id"]}, bson.M{"$set": tw}); err != nil {
		fmt.Println("Error inserting tweet ", tw["id"], "into mongo db.", err)
	} else {
		fmt.Println("Upserted tweet", tw["id"], "to mongo db.")
	}
}

func (mgoC mongoComm) Tweet(id int64) (map[string]interface{}, error) {
	tweetColl := mgoC.session.DB("clayDB").C("Tweet")

	var result map[string]interface{}
	if err := tweetColl.Find(bson.M{"id": id}).One(&result); err != nil {
		if err.Error() == "not found" {
			return map[string]interface{}{}, errors.New("not found")
		}
		fmt.Println("Error finding tweet", id, ".", err)
		return map[string]interface{}{}, errors.New(err.Error())
	}
	return result, nil
}

func (mgoC mongoComm) SaveUser(user map[string]interface{}) {
	userColl := mgoC.session.DB("clayDB").C("User")

	if _, err := userColl.Upsert(bson.M{"id": user["id"].(int64)}, bson.M{"$set": user}); err != nil {
		fmt.Println("Error inserting user ", user["id"], "into mongo db.", err)
	} else {
		fmt.Println("Added user", user["id"], "to mongo db.")
	}
}

func (mgoC mongoComm) User(id int64) (map[string]interface{}, error) {
	userColl := mgoC.session.DB("clayDB").C("User")

	var result map[string]interface{}
	if err := userColl.Find(bson.M{"id": id}).One(&result); err != nil {
		if err.Error() == "not found" {
			return map[string]interface{}{}, errors.New("not found")
		}
		fmt.Println("Error finding user", id, ".", err)
		return map[string]interface{}{}, errors.New(err.Error())
	}
	return result, nil
}

func (mgoC mongoComm) SaveEvent(event map[string]interface{}) {
	eventColl := mgoC.session.DB("clayDB").C("Event")

	if err := eventColl.Insert(event); err != nil {
		fmt.Println("Error inserting event", event, "into mongo db.", err)
	} else {
		fmt.Println("Added event", event["type"], "for", event["id"], "to mongo db.")
	}
}

func (mgoC mongoComm) Event(id int64) (map[string]interface{}, error) {
	eventColl := mgoC.session.DB("clayDB").C("Event")

	var result map[string]interface{}
	if err := eventColl.Find(bson.M{"id": id}).One(&result); err != nil {
		if err.Error() == "not found" {
			return map[string]interface{}{}, errors.New("not found")
		}
		fmt.Println("Error finding event", id, ".", err)
		return map[string]interface{}{}, errors.New(err.Error())
	}
	return result, nil
}

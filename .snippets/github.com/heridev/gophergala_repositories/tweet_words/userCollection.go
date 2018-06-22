package tweet_words

import (
		"gopkg.in/mgo.v2"
		"gopkg.in/mgo.v2/bson"
		"fmt"
		"log"
		)

type User struct {
        Token string
        Secret string
        Keywords []string
}

var GUser User

//var KeywordsArray = []string{"test"}
var KeywordsArray = make(map[string][]string)

func StoreUser(user User) (ret bool) {
		mgoSession, err := mgo.Dial("mongodb://54.188.201.254")
		if err != nil {
			panic(err)
		}
		mgoSession.SetMode(mgo.Monotonic, true)
		defer mgoSession.Close()
		newSes := mgoSession.Copy()
		defer newSes.Close()
		db := newSes.DB("test")
		r, _ := db.CollectionNames()
		fmt.Println(r)
		col := db.C("User")
		if col == nil {
			panic("unable to get collection")
		}
		result := User{}
		fmt.Println(bson.M{"token": user.Token})
		err = col.Find(bson.M{"token": user.Token}).One(&result)
		if result.Token == "" {
			dummy := []string{"bjp"}
			err = col.Insert(&User{user.Token, user.Secret, dummy})
			if err != nil {
				panic(err)
			}
		}
		err = col.Find(bson.M{"token": user.Token}).One(&result)
		if err != nil {
            log.Fatal(err)
        }
	return
}

func StoreKeywords(data string) (ret bool) {
	mgoSession, err := mgo.Dial("mongodb://54.188.201.254")
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Monotonic, true)
	defer mgoSession.Close()
	newSes := mgoSession.Copy()
	defer newSes.Close()
	col := newSes.DB("test").C("User")
	if col == nil {
		panic("unable to get collection")
	}
	result := User{}
	err = col.Find(bson.M{"token": GUser.Token}).One(&result)
	if err != nil {
        log.Fatal(err)
    }
    n := len(result.Keywords)
    words := make([]string, n+1);
    copy(words, result.Keywords[0:])
    words[n] = data
    _, err = col.Upsert(bson.M{"token": GUser.Token, "secret" : GUser.Secret}, bson.M{"$set" : bson.M{"keywords": words}})
    if err != nil {
    	fmt.Println(err)
    }
    result1 := User{}
	err = col.Find(bson.M{"token": GUser.Token}).One(&result1)
	if err != nil {
        log.Fatal(err)
    }
	return
}

func GetKeywords() (ret []string) {
	mgoSession, err := mgo.Dial("mongodb://54.188.201.254")
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Monotonic, true)
	defer mgoSession.Close()
	newSes := mgoSession.Copy()
	defer newSes.Close()
	col := newSes.DB("test").C("User")
	if col == nil {
		panic("unable to get collection")
	}
	result := User{}
	err = col.Find(bson.M{"token": GUser.Token}).One(&result)
	if err != nil {
        log.Fatal(err)
    }
    return result.Keywords
}

func GetTweets(keyValue string) (retValue []TweetStore){
		mgoSession, err := mgo.Dial("mongodb://54.188.201.254")
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Monotonic, true)
	defer mgoSession.Close()
	newSes := mgoSession.Copy()
	defer newSes.Close()
	col := newSes.DB("test").C(keyValue)
	if col == nil {
		panic("unable to get collection")
	}
	result := []TweetStore{}
	err = col.Find(bson.M{}).All(&result)
	if err != nil {
        log.Fatal(err)
    }
    return result
}
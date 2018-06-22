package main

import (
	"github.com/gophergala/opentweet/protocol"
	"github.com/gophergala/opentweet/database"
	"github.com/gophergala/opentweet/rest"
	"log"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("%v", err)
	}
	go serveTweets(db)
	serveRest(db)
}

func serveRest(db database.DB) {
	log.Printf("serving rest?")
	server := rest.NewServer()
	err := server.RegisterUserCB(db.RegisterUser)
	err = server.RegisterTweetCB(db.PostTweet)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("%v", err)
	}
}

func serveTweets(db database.DB) {
	log.Printf("serving tweets?")
	server := protocol.NewServer()
	err := server.Register(db.GetTweets)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("%v", err)
	}
}

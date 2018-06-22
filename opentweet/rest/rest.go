package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Server struct {
	port  int
	user  UserHandler
	tweet TweetHandler
}

type UserHandler struct {
	register UserCallback
}

type TweetHandler struct {
	register TweetCallback
}

type UserCallback func(name, password string) error

type TweetCallback func(user, password, tweet string) error

var port = 8080

func NewServer() Server {
	var server Server
	server.port = port
	return server
}

func (server *Server) RegisterUserCB(fn UserCallback) error {
	server.user.register = fn
	return nil
}

func (server *Server) RegisterTweetCB(fn TweetCallback) error {
	server.tweet.register = fn
	return nil
}

func (handler TweetHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	type tweetpost struct {
		Tweet string `json:"tweet"`
	}
	var post tweetpost

	if req.Method != "POST" {
		resp.WriteHeader(http.StatusBadRequest)
		io.WriteString(resp, "Error not an action\n")
		log.Printf("Request received on /tweet that was not POST")
		return
	}
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&post)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		io.WriteString(resp, fmt.Sprintf("Error could not decode posted data: %v\n", err))
		log.Printf("Posted data on /tweet could not be decoded: %v", err)
		return
	}
	if post.Tweet == "" {
		resp.WriteHeader(http.StatusBadRequest)
		io.WriteString(resp, "Tweet is blank\n")
		log.Printf("Posted data on /tweet, tweet is blank.")
		return
	}
	user, pass, ok := req.BasicAuth()
	if ok != true {
		resp.WriteHeader(http.StatusBadRequest)
		io.WriteString(resp, "No basic auth sent\n")
		log.Printf("Posted data on /tweet, no basic auth sent.")
		return
	}
	err = handler.register(user, pass, post.Tweet)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		io.WriteString(resp, fmt.Sprintf("Error could not post tweet: %v\n", err))
		log.Printf("Tweet could not be posted: %v", err)
		return
	}
}

func (handler UserHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	type userpost struct {
		UserName string `json:"user"`
		Password string `json:"password"`
	}
	var post userpost

	if req.Method != "POST" {
		resp.WriteHeader(http.StatusBadRequest)
		io.WriteString(resp, "Error not an action\n")
		log.Printf("Request received on /users that was not POST")
		return
	}
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&post)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		io.WriteString(resp, fmt.Sprintf("Error could not decode posted data: %v\n", err))
		log.Printf("Posted data on /users could not be decoded: %v", err)
		return
	}
	if post.UserName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		io.WriteString(resp, "User is blank\n")
		log.Printf("Posted data on /users, user is blank.")
		return
	}
	if post.Password == "" {
		resp.WriteHeader(http.StatusBadRequest)
		io.WriteString(resp, "Password is blank\n")
		log.Printf("Posted data on /users, password is blank.")
		return
	}
	err = handler.register(post.UserName, post.Password)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		io.WriteString(resp, fmt.Sprintf("Error could not register user: %v\n", err))
		log.Printf("New user could not be registered: %v", err)
		return
	}
}

func (server *Server) ListenAndServe() error {

	http.Handle("/users", server.user)
	http.Handle("/tweet", server.tweet)

	log.Printf("Listening on %v for rest", server.port)
	serve := fmt.Sprintf(":%v", server.port)
	err := http.ListenAndServe(serve, nil)
	if err != nil {
		return err
	}
	return nil
}

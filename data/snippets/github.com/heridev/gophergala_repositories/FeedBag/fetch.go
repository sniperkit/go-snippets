package main

import (
	"fmt"
	"strconv"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

func StartExistingUsers(c chan<- []Activity) error {
	fmt.Println("Starting go routines")
	u := UserList{}

	err := u.List()
	if err != nil {
		return err
	}

	for _, user := range u {
		StartUserRoutine(user, c)
	}

	return nil
}

func StartUserRoutine(u User, c chan<- []Activity) {
	go userRoutine(u, c)
}

func userRoutine(u User, c chan<- []Activity) error {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: u.AccessToken},
	}

	client := github.NewClient(t.Client())

	//List Options Page, PerPage
	opts := github.ListOptions{1, 50}

	for {
		events, resp, err := client.Activity.ListEventsPerformedByUser(u.Username, false, &opts)
		if err != nil {
			panic(err)
		}

		activityPayload, err := ProcessPayload(events, u)
		activities := ActivityParser(activityPayload)

		c <- activities

		// Wait as long as the X-Poll-Interval header says to
		interval, err := strconv.ParseInt(resp.Header["X-Poll-Interval"][0], 10, 8)
		if err != nil {
			// if strconv failed for whatever reason, use the default X-Poll-Interval value of 60
			time.Sleep(60 * time.Second)
		} else {
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}

	panic("Shouldn't be here")
}

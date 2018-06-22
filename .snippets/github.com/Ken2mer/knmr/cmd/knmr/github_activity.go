package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

func getRepositorySubscription(ctx context.Context, client *github.Client) (*github.Subscription, error) {
	subscription, _, err := client.Activity.GetRepositorySubscription(ctx, username, reponame)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func getEvents(ctx context.Context, client *github.Client) ([]*github.Event, error) {
	events, _, err := client.Activity.ListEventsReceivedByUser(ctx, username, true, nil)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func dumpEvents(events []*github.Event) {
	for _, event := range events {
		fmt.Printf("%s\n", event.GetType())
	}
}

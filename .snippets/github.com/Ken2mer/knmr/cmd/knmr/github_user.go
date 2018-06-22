package main

import (
	"context"
	"fmt"

	"github.com/Ken2mer/knmr/logger"
	"github.com/google/go-github/github"
)

func getFollowingUsers(ctx context.Context, client *github.Client) ([]*github.User, error) {
	users, _, err := client.Users.ListFollowing(ctx, "", nil)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return nil, err
	}
	return users, err
}

func dumpFollowUsers(users []*github.User) {
	for i, user := range users {
		logger.Debugf("%d\n", i)
		fmt.Printf("%s\n", user.GetLogin())
	}
}

func getUser(ctx context.Context, client *github.Client) (*github.User, error) {
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return nil, err
	}
	return user, err
}

func dumpUser(user *github.User) {
	logger.Debugf("\n\n%v\n\n", github.Stringify(user))
}

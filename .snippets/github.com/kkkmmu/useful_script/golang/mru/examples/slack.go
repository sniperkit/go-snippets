package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New("xoxb-212862753029-1kYmZdKDBMfyg6obGq2U9Rnv")
	groups, err := api.GetGroups(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, group := range groups {
		fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
	}

	user, err := api.GetUserInfo("kkkmmu")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
}

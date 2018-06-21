package main

import (
	"fmt"
	"github.com/nlopes/slack"
)

var token = ""
var kkkmmu = ""

func main() {
	//api := slack.New(kkkmmu) //Here use your own token
	api := slack.New(token)
	groups, err := api.GetGroups(false)
	if err != nil {
		panic(err)
	}

	auth, err := api.AuthTest() //The reponse will include your user information
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", auth)

	for _, group := range groups {
		fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
	}

	channels, err := api.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	for _, channel := range channels {
		fmt.Printf("Channel: %s: ID: %s\n", channel.ID, channel.Name)
	}

	params := slack.PostMessageParameters{
		Username:  "ats",
		AsUser:    true,
		Markdown:  true,
		LinkNames: 1,
	}

	/*
			attachment := slack.Attachment{
				Pretext: "some pretext",
				Text:    "*bold* `code` _italic_ ~strike~",
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "a",
						Value: "no",
					},
				},
			}

		params.Attachments = []slack.Attachment{attachment}
	*/
	channelID, timestamp, err := api.PostMessage("G68RDPX8D", "@kkkmmu  *bold* `code` _italic_ ~strike~", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {

		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}

}

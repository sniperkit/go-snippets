package bot

import (
	"encoding/base64"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"regexp"
	"strings"
)

type Bot struct {
	client   *slack.Client
	user     string
	userID   string
	channels map[string]string
	groups   map[string]string
	teamID   string
	team     string
	rtm      *slack.RTM
}

var ATB Bot

var GetAtUserIDFromMessage = regexp.MustCompile(`^\<@(?P<id>[[:word:]]+)\> `)
var ABCDEFG = "eG94Yi0yMTI4NjI3NTMwMjktRFpORUlZTlZlSDl6YmY1d1FtWHdVSXM2"

func (b *Bot) PostMessage(message string) {
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
	_, _, err := b.client.PostMessage("G68RDPX8D", "@kkkmmu *"+message+"*", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func (b *Bot) GetAllChannels() map[string]string {
	return b.channels
}

func (b *Bot) GetAllGroups() map[string]string {
	return b.groups
}

func (b *Bot) GetChannelID(name string) (string, error) {
	if c, ok := b.channels[name]; ok {
		return c, nil
	}

	return "", fmt.Errorf("Channel %s is not exist", name)
}

func (b *Bot) GetGroupID(name string) (string, error) {
	if g, ok := b.groups[name]; ok {
		return g, nil
	}

	return "", fmt.Errorf("Channel %s is not exist", name)
}

func (b *Bot) EventHandler() {
	b.rtm = b.client.NewRTM()
	go b.rtm.ManageConnection()

	for msg := range b.rtm.IncomingEvents {

		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			fmt.Printf("Received Hello: %#v\n", ev)

		case *slack.ConnectedEvent:
			b.rtm.SendMessage(b.rtm.NewOutgoingMessage("Hello Eveyone, I'm Back!", "G68RDPX8D"))

		case *slack.MessageEvent:
			b.HandleMessageEvent(ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %#v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %#v\n", ev.Value)

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

func GetAtMessage(msg *slack.MessageEvent) (string, bool) {
	at := GetAtUserIDFromMessage.FindStringSubmatch(msg.Text)
	if len(at) != 2 {
		return "", false
	}

	fmt.Println(strings.TrimLeft(msg.Text, "<@"+at[1]+">"))
	return strings.TrimLeft(msg.Text, "<@"+at[1]+">"), true
}

func (b *Bot) HandleMessageEvent(msg *slack.MessageEvent) {
	context, ok := GetAtMessage(msg)
	if !ok {
		return
	}

	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(context, "G68RDPX8D"))
	fmt.Println(context)
}

func init() {
	token, err := base64.StdEncoding.DecodeString("eG94Yi0yMTI4NjI3NTMwMjktRFpORUlZTlZlSDl6YmY1d1FtWHdVSXM2")
	if err != nil {
		log.Println("Cannot Get Token:", err)
		return
	}

	ATB.client = slack.New(string(token))
	ATB.groups = make(map[string]string, 1)
	ATB.channels = make(map[string]string, 1)

	auth, err := ATB.client.AuthTest() //The reponse will include your user information
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	ATB.userID = auth.UserID
	ATB.user = auth.User
	ATB.team = auth.Team
	ATB.teamID = auth.TeamID

	channels, err := ATB.client.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	for _, channel := range channels {
		ATB.channels[channel.Name] = channel.ID
	}

	groups, err := ATB.client.GetGroups(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	for _, group := range groups {
		ATB.groups[group.Name] = group.ID
	}

	go ATB.EventHandler()
}

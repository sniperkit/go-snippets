package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/urfave/cli"
)

type lbClient struct {
	bot *linebot.Client
}

var linebotCommand = cli.Command{
	Name:   "linebot",
	Usage:  "linebot",
	Action: linebotCmd,
}

func linebotCmd(ctx *cli.Context) error {
	c, err := newLinebotClient()
	if err != nil {
		return err
	}
	return c.pushMessage()
}

func newLinebotClient() (lbClient, error) {
	// var channelSecret string = "XXXXXX"
	// var channelToken  string = "XXXXXX"
	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return lbClient{}, err
	}
	return lbClient{bot: bot}, nil
}

// cf. https://github.com/line/line-bot-sdk-go#create-message
func (c *lbClient) pushMessage() error {
	var messages []linebot.Message

	leftBtn := linebot.NewMessageTemplateAction("left", "left clicked")
	rightBtn := linebot.NewMessageTemplateAction("right", "right clicked")
	template := linebot.NewConfirmTemplate("Hello World", leftBtn, rightBtn)
	message := linebot.NewTemplateMessage("Sorry :(, please update your app.", template)

	messages = append(messages, message)

	// var userID string = "XXXXXX"
	bot := c.bot
	_, err := bot.PushMessage(userID, messages...).Do()
	if err != nil {
		return err
	}
	return nil
}

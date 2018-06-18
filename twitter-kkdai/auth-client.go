package main

import (
	"fmt"

	"github.com/kkdai/twitter"
)

const (
	//Get consumer key and secret from https://dev.twitter.com/apps/new
	ConsumerKey    string = "YnzeXnXA1vMXMsIrEyq2oNmUT"
	ConsumerSecret string = "EGpXEIBDbRUX98knmf2HHHYTOWXlhvprntPAqmAWGjvBC7e68g"
)

func main() {
	twitterClient := twitter.NewDesktopClient(ConsumerKey, ConsumerSecret)

	//Show a UI to display URL.
	//Please go to this URL to get code to continue
	twitterClient.DoAuth()

	//Get timeline only latest one
	timeline, byteData, err := twitterClient.QueryTimeLine(1)

	if err == nil {
		fmt.Println("timeline struct=", timeline, " byteData=", string(byteData))
	}
}

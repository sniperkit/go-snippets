package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/crowdmob/goamz/sqs"
	"github.com/joho/godotenv"
)

// start of script tag
const ScriptStart = "tell application \"Spotify\" to "

var commands = map[string]string{
	"state":      "player state",
	"play":       "play",
	"pause":      "pause",
	"duration":   "duration of current track",
	"name":       "name of current track",
	"album":      "album of current track",
	"id":         "id of current track",
	"artwork":    "artwork of current track",
	"vol_loud":   "set sound volume to 100",
	"vol_soft":   "set sound volume to 20",
	"vol_norm":   "set sound volume to 50",
	"set_volume": "set sound volume to ", //requires parameter
	"play_track": "play track ",          //requires parameter
	"position":   "player position",
}

func callSpotify(command string, param string) string {
	fullcmd := ScriptStart + commands[command] + param

	out, err := exec.Command("/usr/bin/osascript", "-e", fullcmd).Output()
	if err != nil {
		log.Fatal(err)
		log.Fatal(out)
	}
	return strings.TrimSpace(string(out))
}

//command line processing
// func main() {
//
// 	var cmd = flag.String("o", "pause", "Enter the command for spotify")
// 	flag.Parse()
//
// 	command := callSpotify(*cmd, "\"spotify:track:7p54iuWHqvdeN224OglZ9t\"")
//
// 	if command == "" {
// 		fmt.Println("exiting...")
// 	}
// 	fmt.Println(command)
// 	out, err := exec.Command("/usr/bin/osascript", "-e", command).Output()
// 	if err != nil {
// 		log.Fatal(err)
// 		log.Fatal(out)
// 	}
// done := make(chan bool)
//
// fmt.Println(getTimeLeft())
// go pollSystem()
// <-done
// }

func main() {
	log.Println("Starting sqs processor")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c.AWSAccess = os.Getenv("AWS_ACCESS")
	c.AWSSecret = os.Getenv("AWS_SECRET")
	done := make(chan bool)
	messageQueue := make(chan *sqs.Message)

	//notification queue
	s, err := sqs.NewFrom(c.AWSAccess, c.AWSSecret, "us-east-1")
	if err != nil {
		log.Panic(err)
	}
	q, err := s.GetQueue("spotify-ofp-notification")
	if err != nil {
		log.Panic(err)
	}

	go listenOnQueue("spotify-ofp", messageQueue)
	go processQueue(messageQueue)
	go polling(q)
	<-done
}

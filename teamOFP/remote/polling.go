package main

import (
	"log"
	"strconv"
	"time"

	"github.com/crowdmob/goamz/sqs"
)

// polling sleep time
const sleepTime = time.Second / 2

func polling(queue *sqs.Queue) {
	playerState := getPlayerState()
	track := getCurrentTrackID()
	timeLeft := int(getTimeLeft())
	getNextSong := true
	//inner for loop variables
	var (
		currentPlayerState string
		currentTimeLeft    int
		currentTrack       string
	)
	// log.Println("starting player state: ", playerState)
	for {
		time.Sleep(sleepTime)
		//check player state
		currentPlayerState = getPlayerState()
		currentTimeLeft = int(getTimeLeft())
		currentTrack = getCurrentTrackID()
		// log.Println(currentTrack)
		if playerState != currentPlayerState {
			message := NotificationMessage{"player_" + currentPlayerState, "", currentTrack}
			err := pushMessage(queue, message)
			if err != nil {
				log.Println(err)
			}
			playerState = currentPlayerState
			// log.Println("player state changed: ", currentPlayerState)
		}
		if currentTrack != track {
			if !getNextSong {
				pushMessage(queue, NotificationMessage{"track_ending", track, currentTrack})
				pushMessage(queue, NotificationMessage{"track_start", nextTrack, nextTrack})
				getNextSong = true
				nextTrack := getNextTrack()
				if nextTrack != "" {
					setCurrentTrack(nextTrack)
					track = nextTrack
				} else {
					track = currentTrack
				}
				setNextTrack("")
			}
		}
		//check player duration - is track over
		if currentTimeLeft != timeLeft {
			// log.Println("New Time : ", currentTimeLeft)
			timeLeft = currentTimeLeft
			message := NotificationMessage{"time_left", strconv.Itoa(timeLeft), currentTrack}
			pushMessage(queue, message)
			if timeLeft < 30 && getNextSong { //lock out period
				getNextSong = false
				message := NotificationMessage{"get_next_track", track, currentTrack}
				pushMessage(queue, message)
			}
		}

	}
}

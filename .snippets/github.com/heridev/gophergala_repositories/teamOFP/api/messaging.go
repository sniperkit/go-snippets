package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/crowdmob/goamz/sqs"
)

type NotificationMessage struct {
	Event string `json:"event"`
	Value string `json:"values"`
	Track string `json:"track,omitempty"`
}

func listenOnQueue(q *sqs.Queue, ch chan *sqs.Message) {

	for {
		resp, err := q.ReceiveMessage(1)
		if err != nil {
			log.Panic(err)
		}

		for _, m := range resp.Messages {
			ch <- &m
			q.DeleteMessage(&m)
		}
	}

}

func processQueue(ch chan *sqs.Message) {
	for m := range ch {

		n := &NotificationMessage{}

		if err := json.Unmarshal([]byte(m.Body), &n); err != nil {
			log.Panic(err)
		}

		switch n.Event {
		case "track_end":
			log.Println("Song Ended")
			//queueNextTrack()

		case "get_next_track":
			log.Println("Getting ready to play next track")
			peekNextTrack()

		case "track_start":
			log.Println("Track Started: ", n.Value)
			if n.Track != "" {
				queueNextTrack()
				log.Println("n.Track: ", n.Track)
				splitString := strings.Split(n.Track, ":")
				if len(splitString) > 1 {
					n.Track = splitString[2]
				}
				updateNowPlayingTrack(n.Track)
			}

		case "player_paused":
			log.Println("Player Paused")

		case "player_playing":
			log.Println("Player Playing")

		case "player_stopped":
			log.Println("Player Stopped")

		case "time_left":
			time, _ := strconv.Atoi(n.Value)
			updateNowPlayingTime(time)
			if n.Track != "" {
				log.Println("Time Left: ", n)
				n.Track = strings.Split(n.Track, ":")[2]
				updateNowPlayingTrack(n.Track)
			}

		}
	}
}

func pushMessage(q *sqs.Queue, message interface{}) error {

	j, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = q.SendMessage(string(j))
	if err != nil {
		return err
	}

	return nil
}

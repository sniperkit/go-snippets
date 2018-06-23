package main

import (
	"encoding/json"
	"log"

	"github.com/crowdmob/goamz/sqs"
)

//NotificationMessage simple struct for sending data
type NotificationMessage struct {
	Event string `json:"event"`
	Value string `json:"values"`
	Track string `json:"track"`
}
type context struct {
	AWSAccess string
	AWSSecret string
}

var c = &context{}

func listenOnQueue(queue string, ch chan *sqs.Message) {

	// Setup Queue
	s, err := sqs.NewFrom(c.AWSAccess, c.AWSSecret, "us-east-1")
	if err != nil {
		log.Panic(err)
	}
	q, err := s.GetQueue(queue)
	if err != nil {
		log.Panic(err)
	}

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
		// log.Println("Processing Message: ", m)
		messagebody := map[string]interface{}{}
		err := json.Unmarshal([]byte(m.Body), &messagebody)
		if err != nil {
			log.Panic("the unmarshall plan")
		}
		switch messagebody["command"] {
		case "play_track":
			if str, ok := messagebody["param"].(string); ok {
				setNextTrack("spotify:track:" + str)
			} else {
				log.Panic("was unable to set current track")
			}
			// case "skip_track"
			// case "other command"
		} //end of switch

	}
}

func pushMessage(q *sqs.Queue, message interface{}) error {
	// log.Println("message: ", message)
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

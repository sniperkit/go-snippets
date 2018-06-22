package handler

import (
	"github.com/asim/go-micro/client"

	attendeesproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/attendees"
	readproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/read"
	readuserproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/readuser"
)

func readEvent(id string) (*readproto.Response, error) {
	req := client.NewRequest("gevvent-event-service", "Read.Call", &readproto.Request{
		ID: id,
	})
	rsp := &readproto.Response{}

	// Call service
	if err := client.Call(req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func readUserEvent(eventID, userID string) (*readuserproto.Response, error) {
	req := client.NewRequest("gevvent-event-service", "ReadUser.Call", &readuserproto.Request{
		EventID: eventID,
		UserID:  userID,
	})
	rsp := &readuserproto.Response{}

	// Call service
	if err := client.Call(req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func getAttendees(eventID string) (*attendeesproto.Response, error) {
	req := client.NewRequest("gevvent-event-service", "Attendees.Call", &attendeesproto.Request{
		ID: eventID,
	})
	rsp := &attendeesproto.Response{}

	// Call service
	if err := client.Call(req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

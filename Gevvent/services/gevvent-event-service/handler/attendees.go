package handler

import (
	"code.google.com/p/go.net/context"
	log "github.com/cihub/seelog"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/model"
	"github.com/gophergala/Gevvent/services/gevvent-event-service/repo"

	attendeesproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/attendees"
)

type Attendees struct{}

func (e *Attendees) Call(ctx context.Context, req *attendeesproto.Request, rsp *attendeesproto.Response) error {
	attendees, err := repo.GetAttendees(req.ID)
	if err != nil {
		return err
	}

	for _, attendee := range attendees {
		user, err := readUser(attendee.UserID)
		if err != nil {
			log.Warnf("Error finding user %s, %s", attendee.UserID, err)
		}

		status := attendeesproto.Status_INVITED
		if attendee.Status == model.Going {
			status = attendeesproto.Status_GOING
		}

		rsp.Users = append(rsp.Users, &attendeesproto.User{
			ID:       user.ID,
			Username: user.Username,
			Status:   status,
		})
	}

	return nil
}

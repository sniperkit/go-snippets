package handler

import (
	"fmt"

	"code.google.com/p/go.net/context"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/model"
	"github.com/gophergala/Gevvent/services/gevvent-event-service/repo"

	rsvpproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/rsvp"
)

type RSVP struct{}

func (e *RSVP) Call(ctx context.Context, req *rsvpproto.Request, rsp *rsvpproto.Response) error {
	event, err := repo.Get(req.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return nil
	}

	// If event is private ensure that the user has been invited
	if event.Private && req.UserID != event.UserID {
		userEvent, err := repo.GetUserEvent(req.EventID, req.UserID)
		if err != nil {
			return err
		}
		if userEvent == nil || !(userEvent.Status == model.Going || userEvent.Status == model.Invited) {
			return fmt.Errorf("Permission denied")
		}
	}

	status := model.NotGoing
	switch req.Answer {
	case rsvpproto.Status_NOT_GOING:
		status = model.NotGoing
	case rsvpproto.Status_GOING:
		status = model.Going
	}

	_, err = repo.UpdateUserEvent(&model.UserEvent{
		EventID: req.EventID,
		UserID:  req.UserID,
		Status:  status,
	})
	if err != nil {
		return err
	}

	return nil
}

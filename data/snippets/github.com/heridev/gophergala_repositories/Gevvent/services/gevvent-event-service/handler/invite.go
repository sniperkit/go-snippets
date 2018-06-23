package handler

import (
	"fmt"

	"code.google.com/p/go.net/context"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/model"
	"github.com/gophergala/Gevvent/services/gevvent-event-service/repo"

	inviteproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/invite"
)

type Invite struct{}

func (e *Invite) Call(ctx context.Context, req *inviteproto.Request, rsp *inviteproto.Response) error {
	event, err := repo.Get(req.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return nil
	}

	if event.Private && event.UserID != req.UserID {
		return fmt.Errorf("Permission denied")
	}

	// Try to find user with the requested username
	user, err := readUsername(req.InvitedUser)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("User not found")
	}

	// Check that the user has not already accepted
	userEvent, err := repo.GetUserEvent(req.EventID, user.ID)
	if err != nil {
		return err
	}
	if userEvent != nil && userEvent.Status == model.Going {
		return nil
	}

	_, err = repo.UpdateUserEvent(&model.UserEvent{
		EventID: req.EventID,
		UserID:  user.ID,
		Status:  model.Invited,
	})
	if err != nil {
		return err
	}

	return nil
}

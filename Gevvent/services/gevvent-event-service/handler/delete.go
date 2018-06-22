package handler

import (
	"fmt"

	"code.google.com/p/go.net/context"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/repo"

	deleteproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/delete"
)

type Delete struct{}

func (e *Delete) Call(ctx context.Context, req *deleteproto.Request, rsp *deleteproto.Response) error {
	event, err := repo.Get(req.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return nil
	}

	if event.UserID != req.UserID {
		return fmt.Errorf("Permission denied")
	}

	err = repo.DeleteEvent(req.EventID)
	if err != nil {
		return err
	}

	return nil
}

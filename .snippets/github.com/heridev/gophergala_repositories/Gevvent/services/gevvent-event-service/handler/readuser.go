package handler

import (
	"code.google.com/p/go.net/context"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/repo"

	readuserproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/readuser"
)

type ReadUser struct{}

func (e *ReadUser) Call(ctx context.Context, req *readuserproto.Request, rsp *readuserproto.Response) error {
	event, err := repo.GetUserEvent(req.EventID, req.UserID)
	if err != nil {
		return err
	}
	if event == nil {
		return nil
	}

	rsp.Status = event.Status.ToProto()

	return nil
}

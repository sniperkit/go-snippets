package handler

import (
	"fmt"

	"code.google.com/p/go.net/context"

	"github.com/gophergala/Gevvent/services/gevvent-user-service/repo"

	readuserproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/readuser"
)

type ReadUser struct{}

func (e *ReadUser) ByID(ctx context.Context, req *readuserproto.Request, rsp *readuserproto.Response) error {
	if req.ID == nil {
		return fmt.Errorf("Must specify a user ID")
	}

	user, err := repo.GetUser(req.GetID())
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("Could not find user %s", req.GetID())
	}

	rsp.ID = user.ID
	rsp.Username = user.Username

	return nil
}

func (e *ReadUser) ByUsername(ctx context.Context, req *readuserproto.Request, rsp *readuserproto.Response) error {
	if req.Username == nil {
		return fmt.Errorf("Must specify a user ID")
	}

	user, err := repo.GetUsername(req.GetUsername())
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("Could not find user %s", req.GetUsername())
	}

	rsp.ID = user.ID
	rsp.Username = user.Username

	return nil
}

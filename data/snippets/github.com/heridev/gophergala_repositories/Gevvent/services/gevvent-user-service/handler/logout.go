package handler

import (
	"code.google.com/p/go.net/context"
	"fmt"

	"github.com/gophergala/Gevvent/services/gevvent-user-service/repo"

	logoutproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/logout"
)

type Logout struct{}

func (e *Logout) Call(ctx context.Context, req *logoutproto.Request, rsp *logoutproto.Response) error {
	session, err := checkToken(req.Token)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("Session has expired")
	}

	err = repo.DeleteSession(session.ID)
	if err != nil {
		return err
	}

	return nil
}

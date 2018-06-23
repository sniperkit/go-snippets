package handler

import (
	"code.google.com/p/go.net/context"
	"fmt"

	"github.com/gophergala/Gevvent/services/gevvent-user-service/model"
	"github.com/gophergala/Gevvent/services/gevvent-user-service/repo"

	registerproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/register"
)

type Register struct{}

func (e *Register) Call(ctx context.Context, req *registerproto.Request, rsp *registerproto.Response) error {
	// Check if user already exists
	if user, err := repo.GetUsername(req.Username); err != nil {
		return err
	} else if user != nil {
		return fmt.Errorf("User already exists")
	}

	user, err := model.NewUser(req.Username, req.Password)
	if err != nil {
		return err
	}

	user, err = repo.CreateUser(user)
	if err != nil {
		return err
	}

	rsp = &registerproto.Response{
		UserID: user.ID,
	}

	return nil
}

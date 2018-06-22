package handler

import (
	"fmt"

	"code.google.com/p/go.net/context"
	"github.com/jmcvetta/randutil"
	"golang.org/x/crypto/bcrypt"

	"github.com/gophergala/Gevvent/services/gevvent-user-service/model"
	"github.com/gophergala/Gevvent/services/gevvent-user-service/repo"

	loginproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/login"
)

type Login struct{}

func (e *Login) Call(ctx context.Context, req *loginproto.Request, rsp *loginproto.Response) error {
	user, err := repo.GetUsername(req.Username)
	if err != nil || user == nil {
		return fmt.Errorf("Incorrect username or password")
	} else if bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(req.Password)) != nil {
		return fmt.Errorf("Incorrect username or password")
	}

	// Create new session and return token
	session, err := newSession(user.ID)
	if err != nil {
		return err
	}

	tokenString, err := session.Token()
	if err != nil {
		return err
	}

	rsp.UserID = session.UserID
	rsp.Token = tokenString

	return nil
}

func newSession(userID string) (*model.Session, error) {
	// Create cookie
	secret, err := randutil.AlphaString(64)
	if err != nil {
		return nil, err
	}

	// Create new token
	session := &model.Session{
		UserID: userID,
		Secret: secret,
	}

	// Save token in DB
	session, err = repo.CreateSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

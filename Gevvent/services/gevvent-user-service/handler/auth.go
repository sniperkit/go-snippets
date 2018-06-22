package handler

import (
	"fmt"

	"code.google.com/p/go.net/context"
	"github.com/asim/go-micro/store"
	log "github.com/cihub/seelog"
	"github.com/dgrijalva/jwt-go"
	"github.com/gophergala/Gevvent/services/gevvent-user-service/model"
	"github.com/gophergala/Gevvent/services/gevvent-user-service/repo"

	authorisedproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/authorised"
)

type Authorised struct{}

func (e *Authorised) Call(ctx context.Context, req *authorisedproto.Request, rsp *authorisedproto.Response) error {
	session, err := checkToken(req.Token)
	if err != nil {
		log.Error(err)
		return err
	}
	if session == nil {
		log.Error("Session has expired")
		return fmt.Errorf("Session has expired")
	}

	tokenString, err := session.Token()
	if err != nil {
		log.Error(err)
		return err
	}

	rsp.UserID = session.UserID
	rsp.Token = tokenString

	return nil
}

func checkToken(tokenString string) (*model.Session, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		item, err := store.Get(fmt.Sprintf("auth/signing_key/%s", token.Header["kid"]))
		if err != nil {
			return nil, err
		}

		return item.Value(), nil
	})
	if err == nil && token.Valid {
		var userID, secret string
		if c, ok := token.Claims["uid"].(string); ok {
			userID = c
		}
		if c, ok := token.Claims["sec"].(string); ok {
			secret = c
		}
		// Session not found in DB. Maybe expired?
		session, err := repo.FindSession(userID, secret)
		if err != nil {
			return nil, err
		}

		return session, nil
	}
	if err != nil {
		log.Error(err)
	}

	return nil, nil
}

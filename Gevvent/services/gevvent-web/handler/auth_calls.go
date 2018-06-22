package handler

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/asim/go-micro/client"

	authorisedproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/authorised"
	loginproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/login"
	logoutproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/logout"
	readuserproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/readuser"
	registerproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/register"
)

func checkAuthorised(token string) (userID string, err error) {
	sreq := client.NewRequest("gevvent-user-service", "Authorised.Call", &authorisedproto.Request{
		Token: token,
	})
	srsp := &authorisedproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		return "", err
	}

	return srsp.UserID, nil
}

func loginUser(username, password string) (token string, err error) {
	sreq := client.NewRequest("gevvent-user-service", "Login.Call", &loginproto.Request{
		Username: username,
		Password: password,
	})
	srsp := &loginproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		return "", err
	}

	return srsp.Token, nil
}

func logoutUser(token string) error {
	sreq := client.NewRequest("gevvent-user-service", "Logout.Call", &logoutproto.Request{
		Token: token,
	})
	srsp := &logoutproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		return err
	}

	return nil
}

func registerUser(username, password string) (userID string, err error) {
	sreq := client.NewRequest("gevvent-user-service", "Register.Call", &registerproto.Request{
		Username: username,
		Password: password,
	})
	srsp := &registerproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		return "", err
	}

	return srsp.UserID, nil
}

func readUser(id string) (user *readuserproto.Response, err error) {
	sreq := client.NewRequest("gevvent-user-service", "ReadUser.ByID", &readuserproto.Request{
		ID: proto.String(id),
	})
	srsp := &readuserproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		return nil, err
	}

	return srsp, nil
}

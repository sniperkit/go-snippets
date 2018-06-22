package handler

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/asim/go-micro/client"

	readuserproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/readuser"
)

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

func readUsername(username string) (user *readuserproto.Response, err error) {
	sreq := client.NewRequest("gevvent-user-service", "ReadUser.ByUsername", &readuserproto.Request{
		Username: proto.String(username),
	})
	srsp := &readuserproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		return nil, err
	}

	return srsp, nil
}

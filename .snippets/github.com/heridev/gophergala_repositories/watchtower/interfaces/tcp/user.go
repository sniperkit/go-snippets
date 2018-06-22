package tcp

import (
	"github.com/gophergala/watchtower/messages"
)

// NewUser creates a new user connected over a TCP/IP socket
func newUser() *user {
	return &user{
		messageQueue: make(chan messages.Message, 10),
	}
}

type user struct {
	id           uint32
	messageQueue chan messages.Message
}

func (u *user) ID() uint32 {
	return u.id
}

func (u *user) Send(m messages.Message, channelID uint32) error {
	u.messageQueue <- m
	return nil
}

func (u *user) SetID(id uint32) {
	u.id = id
}

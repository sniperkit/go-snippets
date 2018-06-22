package users

import (
	"github.com/gophergala/watchtower/messages"
)

// A User represents a user that is currently connected
// and listening to one or more channels. The main thing
// tying users together is their ability to receive messages
// via the Send method
type User interface {
	ID() uint32
	Send(messages.Message, uint32) error
	SetID(uint32)
}

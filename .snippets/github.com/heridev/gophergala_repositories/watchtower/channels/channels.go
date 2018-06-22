package channels

import (
	"errors"
	"sync"

	"github.com/gophergala/watchtower/messages"
	"github.com/gophergala/watchtower/users"
)

const (
	channelMessageBufferSize = 100
)

var (
	channels         = make(map[uint32]*Channel)
	channelEditMutex = &sync.RWMutex{}

	// ErrChannelDoesNotExist is thrown if a user tries to
	// send a message to a channel that does not exist
	ErrChannelDoesNotExist = errors.New("the specified channel does not exist")
)

// Join a channel. Creates the channel if it doesn't exist
// This function is idempotent in the sense that one user
// can not "double-join" a channel so successive calls to
// this function have no effect
func Join(userID, channelID uint32) {
	channelEditMutex.Lock()
	defer channelEditMutex.Unlock()

	// Add the new subscriber to an existing channel
	channel, exists := channels[channelID]
	if exists {
		channel.subscribers[userID] = struct{}{}
		channels[channelID] = channel
		return
	}

	// If the channel doesn't exist, create a new
	// one with the caller as the only subscriber
	c := Channel{
		id:          channelID,
		subscribers: map[uint32]struct{}{userID: struct{}{}},
	}

	channels[channelID] = &c
}

// List returns a thread-safe list of active channels
func List() []uint32 {
	var list []uint32

	channelEditMutex.RLock()
	for channelID := range channels {
		list = append(list, channelID)
	}
	channelEditMutex.RUnlock()

	return list
}

// Send sends a message on a channel. The message can be either
// a BroadcastMessage (in which case it will be broadcasted)
// or a PrivateMessage (in which case it will only be sent to
// it's recipients)
func Send(m messages.Message, channelID uint32) error {
	channelEditMutex.RLock()
	defer channelEditMutex.RUnlock()

	// Check that the channel exists
	channel, exists := channels[channelID]
	if !exists {
		return ErrChannelDoesNotExist
	}

	// Send the message to all subscribers
	subscribers := channel.subscribers
	for subscriberID := range subscribers {
		users.Send(subscriberID, channelID, m)
	}

	return nil
}

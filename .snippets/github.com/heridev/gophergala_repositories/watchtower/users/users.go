// Package users manages the registered users
// TODO: Back this with a Bolt DB
package users

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/gophergala/watchtower/messages"
)

const (
	retryTimes = 1000
)

var (
	users      = make(map[uint32]User)
	usersMutex = &sync.RWMutex{}

	// ErrRegisteringNewUser if a new user can't be registered
	ErrRegisteringNewUser = errors.New("error registering new user")
	// NoSuchUserError is thrown when the caller requests details for a non-existant user
	NoSuchUserError = errors.New("no such user registered")
)

// Register registers and stores a new user, returning his ID
func Register(u User) (uint32, error) {
	var userID uint32

	usersMutex.Lock()
	for i := 0; i < retryTimes; i++ {
		// Try retryTimes times to find a random user id, not already in users
		userID = uint32(rand.Int31n(100000))
		_, alreadyExists := users[userID]
		// 0 is not a valid user ID either
		if userID != 0 && !alreadyExists {
			u.SetID(userID)
			users[userID] = u
			break
		}
	}
	usersMutex.Unlock()

	if userID != 0 {
		return userID, nil
	}

	return userID, ErrRegisteringNewUser
}

// List returns a thread-safe copy of the current user list
func List() map[uint32]struct{} {
	// Return a copy of users
	usersCopy := make(map[uint32]struct{})

	usersMutex.RLock()
	for user := range users {
		usersCopy[user] = struct{}{}
	}
	usersMutex.RUnlock()

	return usersCopy
}

// Send a message to a recipient. Errors out if no such recipient or
// message can not be sent for some reason.
func Send(recipientID, channelID uint32, m messages.Message) error {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	user, exists := users[recipientID]
	if !exists {
		return NoSuchUserError
	}

	return user.Send(m, channelID)
}

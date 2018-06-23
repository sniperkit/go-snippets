package channels

import (
	"testing"

	"github.com/gophergala/watchtower/messages"
)

func TestChannels(t *testing.T) {
	Join(1, 1) // Creates a new channel (1), user 1 joins it
	Join(2, 1) // User 2 also joins channel 1

	l := List()
	if len(l) != 1 || l[0] != 1 {
		t.Error("List() returned an invalid list")
	}

	// Try to send on non-existant channel
	err := Send(&messages.BroadcastMessage{}, 2)
	if err != ErrChannelDoesNotExist {
		t.Error("Send should throw ErrChannelDoesNotExist when sending to non-existant channel")
	}

	// Send on the channel that does exist
	// The error from not being able to send to our
	// user (who does not really exist) will be
	// silently ignored
	err = Send(&messages.BroadcastMessage{}, 1)
	if err != nil {
		t.Errorf("Send should work but err was [%v]", err)
	}
}

package users

import (
	"testing"

	"github.com/gophergala/watchtower/messages"
)

type testUser struct{}

func (t *testUser) ID() uint32 {
	return 0
}

func (t *testUser) Send(messages.Message, uint32) error {
	return nil
}

func (t *testUser) SetID(uint32) {
	return
}

func TestUsers(t *testing.T) {
	u := &testUser{}
	id, err := Register(u)
	if err != nil {
		t.Error(err)
	}

	l := List()
	if len(l) != 1 {
		t.Error("invalid List() return value")
	}

	_, exists := l[id]
	if !exists {
		t.Error("invalid List() content")
	}

	Send(id, 1, &messages.BroadcastMessage{})

	// TODO: Figure out a way to test Send()
}

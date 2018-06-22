package tcp

import (
	"testing"
)

func TestUser(t *testing.T) {
	u := newUser()

	if u.ID() != 0 {
		t.Error("invalid id")
	}

	u.SetID(1)

	if u.ID() != 1 {
		t.Error("invalid id after set")
	}
}

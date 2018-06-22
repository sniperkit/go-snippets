package main

import (
	"testing"
	"testing/quick"
)

// Test to ensure room ids we generate are valid.
func TestGeneratedRoomIdAreValid(t *testing.T) {
	f := func() bool {
		roomId := generateRoomId()
		return validateRoomId(roomId) == nil
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

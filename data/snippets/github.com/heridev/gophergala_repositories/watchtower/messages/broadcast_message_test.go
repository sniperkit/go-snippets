package messages

import (
	"testing"
)

var (
	_ Message = (*BroadcastMessage)(nil)
)

func TestBroadcastMessage(t *testing.T) {
	b := BroadcastMessage{
		sender:  0,
		channel: 1,
		content: "hello",
	}

	b2 := NewBroadcastMessage(0, 1, "hello")

	if b != *b2 {
		t.Error("NewBroadcastMessage created wrong type of struct")
	}

	if b.sender != 0 || b.Sender() != 0 {
		t.Error("sender or Sender() did not match")
	}

	if b.channel != 1 || b.Channel() != 1 {
		t.Error("channel or Channel() did not match")
	}

	if b.Receivers() != nil {
		t.Error("receivers or Receivers() did not match")
	}

	if b.content != "hello" || b.Content() != "hello" {
		t.Error("content or Content() did not match")
	}

	if b.JSON() == "" {
		t.Error("failed to encode as JSON")
	}

	bytes := b.Bytes()
	if len(bytes) != 19 {
		t.Error("failed to encode message to bytes")
	}
}

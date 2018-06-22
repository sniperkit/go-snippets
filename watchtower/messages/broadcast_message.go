package messages

import (
	"encoding/binary"
	"encoding/json"
)

// A BroadcastMessage is broadcasted to one or more channels
type BroadcastMessage struct {
	sender  uint32 `json:"id"`
	channel uint32 `json:"channel"`
	content string `json:"content"`
}

// NewBroadcastMessage creates a new broadcast message
func NewBroadcastMessage(sender, channel uint32, content string) *BroadcastMessage {
	return &BroadcastMessage{
		sender:  sender,
		channel: channel,
		content: content,
	}
}

// Sender returns the message's sender ID
func (b *BroadcastMessage) Sender() uint32 {
	return b.sender
}

// Channel returns the ID of the channel the
// the message was sent on
func (b *BroadcastMessage) Channel() uint32 {
	return b.channel
}

// Content returns the message content
func (b *BroadcastMessage) Content() string {
	return b.content
}

// JSON returns a JSON-encoded version of the message
func (b *BroadcastMessage) JSON() string {
	m := make(map[string]interface{})
	m["sender"] = b.sender
	m["content"] = b.content
	m["channel"] = b.channel
	encoded, _ := json.Marshal(m)

	return string(encoded)
}

// Bytes returns a version of the message fit
// for sending over a TCP/IP or UDP pipe. The
// format is defined in the documentation
func (b *BroadcastMessage) Bytes() []byte {
	bytes := make([]byte, (messageHeaderSize + len([]byte(b.content))))

	binary.LittleEndian.PutUint32(bytes[0:4], b.sender)
	binary.LittleEndian.PutUint32(bytes[4:8], b.channel)
	binary.LittleEndian.PutUint16(bytes[8:10], 0) // Message type - reserved for later use
	binary.LittleEndian.PutUint32(bytes[10:14], uint32(len(b.content)))
	copy(bytes[14:], []byte(b.content))
	return bytes
}

// Receivers returns nil for a Broadcast message
func (b *BroadcastMessage) Receivers() map[uint32]struct{} {
	return nil
}
